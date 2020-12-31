package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestDiff(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	commit1, err := Commit(ctx, dag, "/", "1")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	commit2, err := Commit(ctx, dag, "", "2")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	changes, err := Diff(ctx, dag, commit1, commit2)
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
