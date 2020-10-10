// Package file contains utilities for working with files.
package file

import (
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
)

// WriteEntries writes the directory entries to the local repo.
func WriteEntries(entries files.DirIterator, root string) error {
	if !entries.Next() {
		return entries.Err()
	}

	path := filepath.Join(root, entries.Name())
	if err := files.WriteTo(entries.Node(), path); err != nil {
		return err
	}

	return WriteEntries(entries, root)
}