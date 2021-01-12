package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/spf13/afero"
)

func TestDiff(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	treeA, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitA := data.NewCommit(treeA.Cid(), "a")
	idA, err := data.AddCommit(ctx, dag, commitA)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	treeB, err := Add(ctx, dag, "", nil)
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

	if changes[0].Path != "README.md" {
		t.Fatalf("unexpected change path")
	}

	if changes[0].Type != dagutils.Add {
		t.Fatalf("unexpected change type")
	}
}
