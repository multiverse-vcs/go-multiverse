package unixfs

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	ipld "github.com/ipfs/go-ipld-format"
	ufs "github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
)

// Write writes the contents of node to the path.
func Write(ctx context.Context, dag ipld.DAGService, path string, node ipld.Node) error {
	fsnode, err := ufs.ExtractFSNode(node)
	if err != nil {
		return err
	}

	switch fsnode.Type() {
	case ufs.TFile:
		return writeFile(ctx, dag, path, node)
	case ufs.TDirectory:
		return writeDir(ctx, dag, path, node)
	case ufs.TSymlink:
		return os.Symlink(path, string(fsnode.Data()))
	default:
		return errors.New("invalid file type")
	}
}

// writeSymlink writes the file to the given path.
func writeFile(ctx context.Context, dag ipld.DAGService, path string, node ipld.Node) error {
	reader, err := ufsio.NewDagReader(ctx, node, dag)
	if err != nil {
		return err
	}
	defer reader.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := reader.WriteTo(file); err != nil {
		return err
	}

	return nil
}

// writeSymlink writes the directory entries to the given path.
func writeDir(ctx context.Context, dag ipld.DAGService, path string, node ipld.Node) error {
	dir, err := ufsio.NewDirectoryFromNode(dag, node)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	links, err := dir.Links(ctx)
	for _, link := range links {
		subnode, err := link.GetNode(ctx, dag)
		if err != nil {
			return err
		}

		subpath := filepath.Join(path, link.Name)
		if err := Write(ctx, dag, subpath, subnode); err != nil {
			return err
		}
	}

	return nil
}
