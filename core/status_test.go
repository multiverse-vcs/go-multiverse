package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/spf13/afero"
)

func TestStatusRemove(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	tree, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commit := data.NewCommit(tree.Cid(), "init")
	id, err := data.AddCommit(ctx, dag, commit)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := fs.Remove("README.md"); err != nil {
		t.Fatalf("failed to remove readme file")
	}

	changes, err := Status(ctx, dag, "/", nil, id)
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
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	changes, err := Status(ctx, dag, "", nil, cid.Cid{})
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
