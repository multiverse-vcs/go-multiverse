package file

import (
	"fmt"

	"github.com/ipfs/go-ipfs-files"
)

const (
	// Add represents a new file.
	Add    int = iota
	// Remove represents a deleted file.
	Remove
	// Change represents an edited file.
	Change
)

// FileDiff contains info about a change to a file.
type FileDiff struct {
	// Path is the path of the modified file.
	Path string
	// Before is the file before modification.
	Before files.Node
	// After is the file after modification.
	After files.Node
	// Mod is one of Add, Remove, or Change.
	Mod int
}

// Tree is the internal representation of files for building diffs.
type Tree map[string]files.Node

// NewDiff returns a new diff of the file at the given path.
func NewFileDiff(path string, treeA, treeB Tree) *FileDiff {
	nodeA, okA := treeA[path]
	nodeB, okB := treeB[path]

	mod := -1
	switch {
	case !okA && okB:
		mod = Add
	case okA && !okB:
		mod = Remove
	case okA && okB:
		mod = Change
	}

	return &FileDiff{Before: nodeA, After: nodeB, Path: path, Mod: mod}
}

// Diff returns a list of diffs between two file nodes.
func Diff(nodeA, nodeB files.Node) ([]*FileDiff, error) {
	treeA := make(Tree)
	if err := files.Walk(nodeA, treeA.Walk); err != nil {
		return nil, err
	}

	treeB := make(Tree)
	if err := files.Walk(nodeB, treeB.Walk); err != nil {
		return nil, err
	}

	return treeA.Diff(treeB), nil
}

// String returns a string representation of a Diff.
func (d *FileDiff) String() string {
	switch d.Mod {
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

// Diff returns a list of diffs between two trees.
func (treeA Tree) Diff(treeB Tree) []*FileDiff {
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

	diffs := make([]*FileDiff, len(keys))
	for i, name := range keys {
		diffs[i] = NewFileDiff(name, treeA, treeB)
	}

	return diffs
}

// Walk is a callback for building a tree from files.Walk.
func (tree Tree) Walk(path string, node files.Node) error {
	// ignore all but regular files
	if _, ok := node.(files.File); ok {
		tree[path] = node
	}

	return nil
}
