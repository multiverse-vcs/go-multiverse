package core

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/importer/balanced"
	"github.com/ipfs/go-unixfs/importer/helpers"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

// DefaultChunker is the name of the default chunker algorithm.
const DefaultChunker = "buzhash"

// Add creates a node from the file at path and adds it to the merkle dag.
func Add(ctx context.Context, fs afero.Fs, dag ipld.DAGService, path string, filter *ignore.GitIgnore) (ipld.Node, error) {
	stat, err := fs.Stat(path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsRegular():
		return addFile(ctx, fs, dag, path)
	case mode&os.ModeSymlink != 0:
		return addSymlink(ctx, fs, dag, path)
	case mode.IsDir():
		return addDir(ctx, fs, dag, path, filter)
	default:
		return nil, errors.New("invalid file type")
	}
}

func add(ctx context.Context, dag ipld.DAGService, reader io.Reader) (ipld.Node, error) {
	chunker, err := chunk.FromString(reader, DefaultChunker)
	if err != nil {
		return nil, err
	}

	params := helpers.DagBuilderParams{
		Dagserv:    dag,
		CidBuilder: merkledag.V1CidPrefix(),
		Maxlinks:   helpers.DefaultLinksPerBlock,
	}

	helper, err := params.New(chunker)
	if err != nil {
		return nil, err
	}

	node, err := balanced.Layout(helper)
	if err != nil {
		return nil, err
	}

	return node, dag.Add(ctx, node)
}

func addFile(ctx context.Context, fs afero.Fs, dag ipld.DAGService, path string) (ipld.Node, error) {
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return add(ctx, dag, file)
}

func addSymlink(ctx context.Context, fs afero.Fs, dag ipld.DAGService, path string) (ipld.Node, error) {
	reader, ok := fs.(afero.LinkReader)
	if !ok {
		return nil, errors.New("fs does not support symlinks")
	}

	target, err := reader.ReadlinkIfPossible(path)
	if err != nil {
		return nil, err
	}

	data, err := unixfs.SymlinkData(target)
	if err != nil {
		return nil, err
	}

	node := merkledag.NodeWithData(data)
	return node, dag.Add(ctx, node)
}

func addDir(ctx context.Context, fs afero.Fs, dag ipld.DAGService, path string, filter *ignore.GitIgnore) (ipld.Node, error) {
	entries, err := afero.ReadDir(fs, path)
	if err != nil {
		return nil, err
	}

	dir := ufsio.NewDirectory(dag)
	for _, info := range entries {
		subpath := filepath.Join(path, info.Name())
		if filter != nil && filter.MatchesPath(subpath) {
			continue
		}

		subnode, err := Add(ctx, fs, dag, subpath, filter)
		if err != nil {
			return nil, err
		}

		if err := dir.AddChild(ctx, info.Name(), subnode); err != nil {
			return nil, err
		}
	}

	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}

	return node, dag.Add(ctx, node)
}
