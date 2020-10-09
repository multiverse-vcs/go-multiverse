package core

import (
	"context"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gookit/color"
	"github.com/ipfs/go-ipfs-files"
)

// WriteTreeEntries writes the tree entries to the local repo.
func WriteTreeEntries(entries files.DirIterator, root string) error {
	if !entries.Next() {
		return entries.Err()
	}

	path := filepath.Join(root, entries.Name())
	if err := files.WriteTo(entries.Node(), path); err != nil {
		return err
	}

	return WriteTreeEntries(entries, root)
}

// MapTreeEntries flattens and maps entries by path.
func MapTreeEntries(node files.Node, paths map[string]files.Node) error {
	cb := func(p string, n files.Node) error {
		if strings.Compare(p, "") == 0 {
			return nil
		}

		paths[p] = n
		return nil
	}

	return files.Walk(node, cb)
}

// DiffTrees prints the differences between two trees.
func DiffTrees(ctx context.Context, nodeA, nodeB files.Node) error {
	treeA := make(map[string]files.Node)
	if err := MapTreeEntries(nodeA, treeA); err != nil {
		return err
	}

	treeB := make(map[string]files.Node)
	if err := MapTreeEntries(nodeB, treeB); err != nil {
		return err
	}

	paths := make(map[string]bool)
	for name := range treeA {
		paths[name] = true
	}

	for name := range treeB {
		paths[name] = true
	}

	keys := make([]string, len(paths))
	for name := range paths {
		keys = append(keys, name)
	}

	sort.Strings(keys)
	for _, name := range keys {
		_, okA := treeA[name]
		_, okB := treeB[name]

		switch {
		case okA && !okB:
			color.Red.Printf("\tdeleted: %s\n", name)
		case !okA && okB:
			color.Green.Printf("\tnew file: %s\n", name)
		case okA && okB:
			color.Yellow.Printf("\tmodified: %s\n", name)
		}
	}

	return nil
}