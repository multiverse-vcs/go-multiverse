package core

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
)

func writeNode(node files.Node, root string) error {
	dir, ok := node.(files.Directory)
	if ok {
		return writeDirectory(dir, root)
	}

	file, ok := node.(files.File)
	if ok {
		return writeFile(file, root)
	}

	return ErrInvalidFile
}

func writeFile(node files.File, root string) error {
	b, err := ioutil.ReadAll(node)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(root, b, 0644)
}

func writeDirectory(node files.Directory, root string) error {
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	return writeEntries(node.Entries(), root)
}

func writeEntries(entries files.DirIterator, root string) error {
	if !entries.Next() {
		return entries.Err()
	}

	path := filepath.Join(root, entries.Name())
	if err := writeNode(entries.Node(), path); err != nil {
		return err
	}

	return writeEntries(entries, root)
}
