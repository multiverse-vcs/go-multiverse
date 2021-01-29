package unixfs

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	ufs "github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
)

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

// addFile creates a dag node from the file at the given path.
func addFile(ctx context.Context, dag ipld.DAGService, path string) (ipld.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Chunk(ctx, dag, file)
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
	if err := dag.Add(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
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

	if err := dag.Add(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}
