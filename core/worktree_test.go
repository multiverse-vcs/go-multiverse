package core

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-ipfs-files"
)

func TestWorktree(t *testing.T) {
	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}

	ignore := filepath.Join(c.root, IgnoreFile)
	if err := ioutil.WriteFile(ignore, []byte("*.exe"), 0644); err != nil {
		t.Fatalf("failed to write ignore file")
	}

	exe := filepath.Join(c.root, "test.exe")
	if err := ioutil.WriteFile(exe, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write exe file")
	}

	tree, err := c.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree")
	}

	node, ok := tree.(files.Directory)
	if !ok {
		t.Fatalf("expected tree to be a directory")
	}

	entries := node.Entries()
	if !entries.Next() {
		t.Fatalf("tree should not be empty")
	}

	if entries.Name() != ".multignore" {
		t.Fatalf("unexpected tree entry: %s", entries.Name())
	}

	if entries.Next() {
		t.Fatalf("unexpected tree entry: %s", entries.Name())
	}
}
