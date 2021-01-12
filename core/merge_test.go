package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/spf13/afero"
)

func TestMergeConflicts(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	treeO, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitO := data.NewCommit(treeO.Cid(), "o")
	o, err := data.AddCommit(ctx, dag, commitO)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello\nfoo\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	treeA, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	commitA := data.NewCommit(treeA.Cid(), "a", o)
	a, err := data.AddCommit(ctx, dag, commitA)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello\nbar\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	treeB, err := Add(ctx, dag, "", nil)
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
		t.Fatalf("failed to merge %s", err)
	}
}
