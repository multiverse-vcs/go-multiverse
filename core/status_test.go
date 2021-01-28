package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

func TestStatusBare(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	changes, err := Status(ctx, dag, "testdata/2", nil, cid.Cid{})
	if err != nil {
		t.Fatalf("failed to get status")
	}

	if len(changes) != 1 {
		t.Fatalf("unexpected changes")
	}

	if changes[0].Path != "README.txt" {
		t.Fatalf("unexpected change path")
	}

	if changes[0].Type != dagutils.Add {
		t.Fatalf("unexpected change type")
	}
}

func TestStatus(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tree, err := unixfs.Add(ctx, dag, "testdata/1", unixfs.Ignore{".gitkeep"})
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commit := data.NewCommit(tree.Cid(), "init")
	id, err := data.AddCommit(ctx, dag, commit)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	changes, err := Status(ctx, dag, "testdata/2", nil, id)
	if err != nil {
		t.Fatalf("failed to get status")
	}

	if len(changes) != 1 {
		t.Fatalf("unexpected changes")
	}

	if changes[0].Path != "README.txt" {
		t.Fatalf("unexpected change path")
	}

	if changes[0].Type != dagutils.Add {
		t.Fatalf("unexpected change type")
	}
}
