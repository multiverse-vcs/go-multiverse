package fs

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	unixfs "github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/io"

	"github.com/multiverse-vcs/go-multiverse/internal/ignore"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
)

// Add creates a node from the file at path and adds it to the merkle dag.
func Add(ctx context.Context, ds ipld.DAGService, path string, filter ignore.Filter) (ipld.Node, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsRegular():
		return addFile(ctx, ds, path)
	case mode&os.ModeSymlink != 0:
		return addSymlink(ctx, ds, path)
	case mode.IsDir():
		return addDir(ctx, ds, path, filter)
	default:
		return nil, errors.New("invalid file type")
	}
}

// addFile creates a dag node from the file at the given path.
func addFile(ctx context.Context, ds ipld.DAGService, path string) (ipld.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return dag.Chunk(ctx, ds, file)
}

// addSymlink creates a dag node from the symlink at the given path.
func addSymlink(ctx context.Context, ds ipld.DAGService, path string) (ipld.Node, error) {
	target, err := os.Readlink(path)
	if err != nil {
		return nil, err
	}

	data, err := unixfs.SymlinkData(target)
	if err != nil {
		return nil, err
	}

	node := merkledag.NodeWithData(data)
	if err := ds.Add(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

// addDir creates a dag node from the directory entries at the given path.
func addDir(ctx context.Context, ds ipld.DAGService, path string, filter ignore.Filter) (ipld.Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	other, err := ignore.Load(path)
	if err != nil {
		return nil, err
	}

	filter = filter.Merge(other)
	ufsdir := io.NewDirectory(ds)

	for _, info := range entries {
		subpath := filepath.Join(path, info.Name())
		if filter.Match(subpath) {
			continue
		}

		subnode, err := Add(ctx, ds, subpath, filter)
		if err != nil {
			return nil, err
		}

		if err := ufsdir.AddChild(ctx, info.Name(), subnode); err != nil {
			return nil, err
		}
	}

	node, err := ufsdir.GetNode()
	if err != nil {
		return nil, err
	}

	if err := ds.Add(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}
