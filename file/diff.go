package file

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/ipfs/go-ipfs-files"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	// Add represents a new file.
	Add int = iota
	// Remove represents a deleted file.
	Remove
	// Change represents an edited file.
	Change
)

// FileDiff contains info about a change to a file.
type FileDiff struct {
	// Mod is one of Add, Remove, or Change.
	Mod int
	// Path is the path of the modified file.
	Path string
	// Before is the file before modification.
	Before files.File
	// After is the file after modification.
	After files.File
}

// Tree is the internal representation of files for building diffs.
type Tree map[string]files.File

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

	return treeA.Diff(treeB)
}

// NewFileDiff creates a file diff from two trees.
func NewFileDiff(treeA, treeB Tree, path string) *FileDiff {
	nodeA, okA := treeA[path]
	nodeB, okB := treeB[path]

	if !okA && okB {
		return &FileDiff{Add, path, nodeA, nodeB}
	}

	if okA && !okB {
		return &FileDiff{Remove, path, nodeA, nodeB}
	}

	return &FileDiff{Change, path, nodeA, nodeB}
}

// Patch returns a string representation of the changes in the diff.
func (d *FileDiff) Patch() (string, error) {
	bytesA, err := ioutil.ReadAll(d.Before)
	if err != nil {
		return "", err
	}

	bytesB, err := ioutil.ReadAll(d.After)
	if err != nil {
		return "", err
	}

	if bytes.Equal(bytesA, bytesB) {
		return "", err
	}

	dmp := diffmatchpatch.New()
	runesA, runesB, lines := dmp.DiffLinesToRunes(string(bytesA), string(bytesB))
	diffs := dmp.DiffMainRunes(runesA, runesB, false)
	diffs = dmp.DiffCharsToLines(diffs, lines)
	return dmp.DiffPrettyText(diffs), nil
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
func (treeA Tree) Diff(treeB Tree) ([]*FileDiff, error) {
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
		diffs[i] = NewFileDiff(treeA, treeB, name)
	}

	return diffs, nil
}

// Walk is a callback for building a tree from files.Walk.
func (tree Tree) Walk(path string, node files.Node) error {
	if fn, ok := node.(files.File); ok {
		tree[path] = fn
	}

	return nil
}
