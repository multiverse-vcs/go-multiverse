package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
	"github.com/ipfs/go-merkledag/dagutils"
)

func TestStatusRemove(t *testing.T) {
	mock := NewMockContext()

	readme := mock.fs.Join(mock.fs.Root(), "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	if _, err := mock.Commit("init"); err != nil {
		t.Fatalf("failed to commit")
	}

	if err := mock.fs.Remove(readme); err != nil {
		t.Fatalf("failed to remove readme file")
	}

	changes, err := mock.Status()
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
	mock := NewMockContext()

	readme := mock.fs.Join(mock.fs.Root(), "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	changes, err := mock.Status()
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
