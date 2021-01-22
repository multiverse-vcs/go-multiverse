package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

func TestMergeConflicts(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	treeO, err := unixfs.Add(ctx, dag, "testdata/2", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitO := data.NewCommit(treeO.Cid(), "o")
	o, err := data.AddCommit(ctx, dag, commitO)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	treeA, err := unixfs.Add(ctx, dag, "testdata/3", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitA := data.NewCommit(treeA.Cid(), "a", o)
	a, err := data.AddCommit(ctx, dag, commitA)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	treeB, err := unixfs.Add(ctx, dag, "testdata/4", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitB := data.NewCommit(treeB.Cid(), "b", o)
	b, err := data.AddCommit(ctx, dag, commitB)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	_, err = Merge(ctx, dag, o, a, b)
	if err != nil {
		t.Fatalf("failed to merge")
	}
}
