package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestStatusRemove(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	head, err := Commit(context.TODO(), fs, dag, "init")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := fs.Remove("README.md"); err != nil {
		t.Fatalf("failed to remove readme file")
	}

	changes, err := Status(context.TODO(), fs, dag, head)
	if err != nil {
		t.Fatalf("failed to get status")
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
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	changes, err := Status(context.TODO(), fs, dag, cid.Cid{})
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
