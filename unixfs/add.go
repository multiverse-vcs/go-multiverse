package unixfs

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	ufs "github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/importer/balanced"
	"github.com/ipfs/go-unixfs/importer/helpers"
	ufsio "github.com/ipfs/go-unixfs/io"
)

// DefaultChunker is the name of the default chunker algorithm.
const DefaultChunker = "buzhash"

// Ignore is used to filter files.
type Ignore []string

// Match returns true if the path matches any ignore rules.
func (i Ignore) Match(path string) bool {
	for _, p := range i {
		base := filepath.Base(path)
		if match, _ := filepath.Match(p, base); match {
			return true
		}

		if match, _ := filepath.Match(p, path); match {
			return true
		}
	}

	return false
}

// Add creates a node from the file at path and adds it to the merkle dag.
func Add(ctx context.Context, dag ipld.DAGService, path string, ignore Ignore) (ipld.Node, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsRegular():
		return addFile(ctx, dag, path)
	case mode&os.ModeSymlink != 0:
		return addSymlink(ctx, dag, path)
	case mode.IsDir():
		return addDir(ctx, dag, path, ignore)
	default:
		return nil, errors.New("invalid file type")
	}
}

// addReader splits the given reader into chunks and arranges them into a dag node.
func addReader(ctx context.Context, dag ipld.DAGService, reader io.Reader) (ipld.Node, error) {
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

// addFile creates a dag node from the file at the given path.
func addFile(ctx context.Context, dag ipld.DAGService, path string) (ipld.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return addReader(ctx, dag, file)
}

// addSymlink creates a dag node from the symlink at the given path.
func addSymlink(ctx context.Context, dag ipld.DAGService, path string) (ipld.Node, error) {
	target, err := os.Readlink(path)
	if err != nil {
		return nil, err
	}

	data, err := ufs.SymlinkData(target)
	if err != nil {
		return nil, err
	}

	node := merkledag.NodeWithData(data)
	return node, dag.Add(ctx, node)
}

// addDir creates a dag node from the directory entries at the given path.
func addDir(ctx context.Context, dag ipld.DAGService, path string, ignore Ignore) (ipld.Node, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	dir := ufsio.NewDirectory(dag)
	for _, info := range entries {
		subpath := filepath.Join(path, info.Name())
		if ignore.Match(subpath) {
			continue
		}

		subnode, err := Add(ctx, dag, subpath, ignore)
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
