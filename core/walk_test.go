package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

func TestWalk(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tree, err := unixfs.Add(ctx, dag, "testdata/1", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitA := data.NewCommit(tree.Cid(), "first")
	idA, err := data.AddCommit(ctx, dag, commitA)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	commitB := data.NewCommit(tree.Cid(), "second", idA)
	idB, err := data.AddCommit(ctx, dag, commitB)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	history, err := Walk(ctx, dag, idB, nil)
	if err != nil {
		t.Fatalf("failed to walk")
	}

	if len(history) != 2 {
		t.Fatalf("cids do not match")
	}

	if !history[idA.KeyString()] {
		t.Errorf("cids do not match")
	}

	if !history[idB.KeyString()] {
		t.Errorf("cids do not match")
	}
}
