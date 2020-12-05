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
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

// DefaultChunker is the name of the default chunker algorithm.
const DefaultChunker = "buzhash"

// Add creates a node from the file at path and adds it to the merkle dag.
func Add(ctx context.Context, store *storage.Store, path string, filter *ignore.GitIgnore) (ipld.Node, error) {
	stat, err := store.Cwd.Stat(path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsRegular():
		return addFile(ctx, store, path)
	case mode&os.ModeSymlink != 0:
		return addSymlink(ctx, store, path)
	case mode.IsDir():
		return addDir(ctx, store, path, filter)
	default:
		return nil, errors.New("invalid file type")
	}
}

func add(ctx context.Context, store *storage.Store, reader io.Reader) (ipld.Node, error) {
	chunker, err := chunk.FromString(reader, DefaultChunker)
	if err != nil {
		return nil, err
	}

	params := helpers.DagBuilderParams{
		Dagserv:    store.Dag,
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

	return node, store.Dag.Add(ctx, node)
}

func addFile(ctx context.Context, store *storage.Store, path string) (ipld.Node, error) {
	file, err := store.Cwd.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return add(ctx, store, file)
}

func addSymlink(ctx context.Context, store *storage.Store, path string) (ipld.Node, error) {
	reader, ok := store.Cwd.(afero.LinkReader)
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
	return node, store.Dag.Add(ctx, node)
}

func addDir(ctx context.Context, store *storage.Store, path string, filter *ignore.GitIgnore) (ipld.Node, error) {
	entries, err := afero.ReadDir(store.Cwd, path)
	if err != nil {
		return nil, err
	}

	dir := ufsio.NewDirectory(store.Dag)
	for _, info := range entries {
		subpath := filepath.Join(path, info.Name())
		if filter != nil && filter.MatchesPath(subpath) {
			continue
		}

		subnode, err := Add(ctx, store, subpath, filter)
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

	return node, store.Dag.Add(ctx, node)
}
