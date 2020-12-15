package core

import (
	"context"
	"errors"
	"path/filepath"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/spf13/afero"
)

// Write writes the contents of node to the path.
func Write(ctx context.Context, fs afero.Fs, dag ipld.DAGService, path string, node ipld.Node) error {
	fsnode, err := unixfs.ExtractFSNode(node)
	if err != nil {
		return err
	}

	switch fsnode.Type() {
	case unixfs.TFile:
		return writeFile(ctx, fs, dag, path, node)
	case unixfs.TDirectory:
		return writeDir(ctx, fs, dag, path, node)
	case unixfs.TSymlink:
		return writeSymlink(fs, path, fsnode.Data())
	default:
		return errors.New("invalid file type")
	}
}

func writeSymlink(fs afero.Fs, path string, target []byte) error {
	linker, ok := fs.(afero.Linker)
	if !ok {
		return errors.New("fs does not support symlinks")
	}

	return linker.SymlinkIfPossible(path, string(target))
}

func writeFile(ctx context.Context, fs afero.Fs, dag ipld.DAGService, path string, node ipld.Node) error {
	reader, err := ufsio.NewDagReader(ctx, node, dag)
	if err != nil {
		return err
	}
	defer reader.Close()

	file, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := reader.WriteTo(file); err != nil {
		return err
	}

	return nil
}

func writeDir(ctx context.Context, fs afero.Fs, dag ipld.DAGService, path string, node ipld.Node) error {
	dir, err := ufsio.NewDirectoryFromNode(dag, node)
	if err != nil {
		return err
	}

	if err := fs.MkdirAll(path, 0755); err != nil {
		return err
	}

	links, err := dir.Links(ctx)
	for _, link := range links {
		subnode, err := link.GetNode(ctx, dag)
		if err != nil {
			return err
		}

		subpath := filepath.Join(path, link.Name)
		if err := Write(ctx, fs, dag, subpath, subnode); err != nil {
			return err
		}
	}

	return nil
}
