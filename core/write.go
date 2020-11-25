package core

import (
	"context"
	"errors"
	"path/filepath"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

// Write writes the contents of node to the path.
func Write(ctx context.Context, store *storage.Store, path string, node ipld.Node) error {
	fsnode, err := unixfs.ExtractFSNode(node)
	if err != nil {
		return err
	}

	switch fsnode.Type() {
	case unixfs.TFile:
		return writeFile(ctx, store, path, node)
	case unixfs.TDirectory:
		return writeDir(ctx, store, path, node)
	case unixfs.TSymlink:
		return writeSymlink(store, path, fsnode.Data())
	default:
		return errors.New("invalid file type")
	}
}

func writeSymlink(store *storage.Store, path string, target []byte) error {
	linker, ok := store.Cwd.(afero.Linker)
	if !ok {
		return errors.New("fs does not support symlinks")
	}

	return linker.SymlinkIfPossible(path, string(target))
}

func writeFile(ctx context.Context, store *storage.Store, path string, node ipld.Node) error {
	reader, err := ufsio.NewDagReader(ctx, node, store.Dag)
	if err != nil {
		return err
	}

	file, err := store.Cwd.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := reader.WriteTo(file); err != nil {
		return err
	}

	return nil
}

func writeDir(ctx context.Context, store *storage.Store, path string, node ipld.Node) error {
	dir, err := ufsio.NewDirectoryFromNode(store.Dag, node)
	if err != nil {
		return err
	}

	if err := store.Cwd.MkdirAll(path, 0755); err != nil {
		return err
	}

	links, err := dir.Links(ctx)
	for _, link := range links {
		subnode, err := link.GetNode(ctx, store.Dag)
		if err != nil {
			return err
		}

		subpath := filepath.Join(path, link.Name)
		if err := Write(ctx, store, subpath, subnode); err != nil {
			return err
		}
	}

	return nil
}
