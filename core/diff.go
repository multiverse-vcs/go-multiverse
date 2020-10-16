package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
)

const (
	// Add represents a new file.
	Add int = iota
	// Remove represents a deleted file.
	Remove
	// Change represents an edited file.
	Change
)

// Diff contains info about a change to a file.
type Diff struct {
	// Type is one of Add, Remove, or Change.
	Type int
	// Path is the path of the modified file.
	Path string
	// Before is the file before modification.
	Before files.File
	// After is the file after modification.
	After files.File
}

type tree map[string]files.File

func (tree tree) walk(path string, node files.Node) error {
	if fn, ok := node.(files.File); ok {
		tree[path] = fn
	}

	return nil
}

func (treeA tree) diff(treeB tree, name string) (*Diff, error) {
	nodeA, okA := treeA[name]
	nodeB, okB := treeB[name]

	if !okA && okB {
		return &Diff{Add, name, nodeA, nodeB}, nil
	}

	if okA && !okB {
		return &Diff{Remove, name, nodeA, nodeB}, nil
	}

	bytesA, err := ioutil.ReadAll(nodeA)
	if err != nil {
		return nil, err
	}

	bytesB, err := ioutil.ReadAll(nodeB)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(bytesA, bytesB) {
		return nil, nil
	}

	return &Diff{Change, name, nodeA, nodeB}, nil
}

// BuildDiffs returns a list of diffs between the files.
func BuildDiffs(nodeA, nodeB files.Node) ([]*Diff, error) {
	treeA := make(tree)
	if err := files.Walk(nodeA, treeA.walk); err != nil {
		return nil, err
	}

	treeB := make(tree)
	if err := files.Walk(nodeB, treeB.walk); err != nil {
		return nil, err
	}

	paths := make(map[string]bool)
	for name := range treeA {
		paths[name] = true
	}

	for name := range treeB {
		paths[name] = true
	}

	keys := make([]string, 0, len(paths))
	for name := range paths {
		keys = append(keys, name)
	}

	diffs := make([]*Diff, 0, len(keys))
	for _, name := range keys {
		diff, err := treeA.diff(treeB, name)
		if err != nil {
			return nil, err
		}

		if diff != nil {
			diffs = append(diffs, diff)
		}
	}

	return diffs, nil
}

// Apply applies the diff to the local repo.
func (d *Diff) Apply(root string) error {
	path := filepath.Join(root, d.Path)
	if d.Type == Remove {
		return os.Remove(path)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(d.After)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)
}

// String returns a string representation of a Diff.
func (d *Diff) String() string {
	switch d.Type {
	case Add:
		return fmt.Sprintf("add:    %s", d.Path)
	case Remove:
		return fmt.Sprintf("remove: %s", d.Path)
	case Change:
		return fmt.Sprintf("change: %s", d.Path)
	default:
		return "unknown diff type"
	}
}
