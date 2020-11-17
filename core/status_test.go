package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestStatusRemove(t *testing.T) {
	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}

	readme := filepath.Join(c.root, "README.md")
	if err := ioutil.WriteFile(readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write readme file")
	}

	if _, err := c.Commit("init"); err != nil {
		t.Fatalf("failed to commit")
	}

	if err := os.Remove(readme); err != nil {
		t.Fatalf("failed to remove readme file")
	}

	changes, err := c.Status()
	if err != nil {
		t.Fatalf("failed to get status: %s", err)
	}

	if len(changes) != 1 {
		t.Fatalf("unexpected changes")
	}

	if changes[0].Path != "README.md" {
		t.Fatalf("unexpected change path")
	}

	if changes[0].Type != dagutils.Remove {
		t.Fatalf("unexpected change type")
	}
}

func TestStatusBare(t *testing.T) {
	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}

	readme := filepath.Join(c.root, "README.md")
	if err := ioutil.WriteFile(readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write readme file")
	}

	changes, err := c.Status()
	if err != nil {
		t.Fatalf("failed to get status")
	}

	if len(changes) != 1 {
		t.Fatalf("unexpected changes")
	}

	if changes[0].Path != "README.md" {
		t.Fatalf("unexpected change path")
	}

	if changes[0].Type != dagutils.Add {
		t.Fatalf("unexpected change type")
	}
}
