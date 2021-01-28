package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

func TestDiff(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	treeA, err := unixfs.Add(ctx, dag, "testdata/1", unixfs.Ignore{".gitkeep"})
	if err != nil {
		t.Fatalf("failed to add tree %s", err)
	}

	commitA := data.NewCommit(treeA.Cid(), "a")
	idA, err := data.AddCommit(ctx, dag, commitA)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	treeB, err := unixfs.Add(ctx, dag, "testdata/2", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitB := data.NewCommit(treeB.Cid(), "b")
	idB, err := data.AddCommit(ctx, dag, commitB)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	changes, err := Diff(ctx, dag, idA, idB)
	if err != nil {
		t.Fatalf("failed to get diff")
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
