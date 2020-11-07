// Package util contains helpers that don't belong elsewhere.
package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
)

// WriteTo writes the files node to the root path.
func WriteTo(node files.Node, root string) error {
	switch node := node.(type) {
	case *files.Symlink:
		return os.Symlink(node.Target, root)
	case files.File:
		b, err := ioutil.ReadAll(node)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(root, b, 0644)
	case files.Directory:
		if err := os.MkdirAll(root, 0777); err != nil {
			return err
		}

		entries := node.Entries()
		for entries.Next() {
			child := filepath.Join(root, entries.Name())
			if err := WriteTo(entries.Node(), child); err != nil {
				return err
			}
		}

		return entries.Err()
	default:
		return fmt.Errorf("invalid file type")
	}
}
