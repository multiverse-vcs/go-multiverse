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

// Tree contains a flattened file tree.
type FileTree map[string]files.File

// Walk is a callback for building a tree from a node.
func (tree FileTree) Walk(path string, node files.Node) error {
	if fn, ok := node.(files.File); ok {
		tree[path] = fn
	}

	return nil
}

// Diff returns a list of diffs between two file nodes.
func Diff(nodeA, nodeB files.Node) ([]*FileDiff, error) {
	treeA := make(FileTree)
	if err := files.Walk(nodeA, treeA.Walk); err != nil {
		return nil, err
	}

	treeB := make(FileTree)
	if err := files.Walk(nodeB, treeB.Walk); err != nil {
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

	diffs := make([]*FileDiff, 0, len(keys))
	for _, name := range keys {
		nodeA, okA := treeA[name]
		nodeB, okB := treeB[name]

		if !okA && okB {
			diffs = append(diffs, &FileDiff{Add, name, nodeA, nodeB})
			continue 
		}

		if okA && !okB {
			diffs = append(diffs, &FileDiff{Remove, name, nodeA, nodeB})
			continue
		}

		bytesA, err := ioutil.ReadAll(nodeA)
		if err != nil {
			return nil, err
		}

		bytesB, err := ioutil.ReadAll(nodeB)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(bytesA, bytesB) {
			diffs = append(diffs, &FileDiff{Change, name, nodeA, nodeB})
		}
	}

	return diffs, nil
}

// BeforeString returns the contents of before as a string.
func (d *FileDiff) BeforeString() (string, error) {
	if d.Before == nil {
		return "", nil
	}

	b, err := ioutil.ReadAll(d.Before)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// AfterString returns the contents of after as a string.
func (d *FileDiff) AfterString() (string, error) {
	if d.After == nil {
		return "", nil
	}

	b, err := ioutil.ReadAll(d.After)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Patch returns a string representation of the changes in the diff.
func (d *FileDiff) Patch() (string, error) {
	before, err := d.BeforeString()
	if err != nil {
		return "", err
	}

	after, err := d.AfterString()
	if err != nil {
		return "", err
	}

	dmp := diffmatchpatch.New()
	runesA, runesB, lines := dmp.DiffLinesToRunes(before, after)
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
