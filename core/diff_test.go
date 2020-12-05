package core

import (
	"context"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
	"testing"
)

func TestDiff(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	commit1, err := Commit(context.TODO(), store, "1")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	commit2, err := Commit(context.TODO(), store, "2")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	changes, err := Diff(context.TODO(), store, commit1, commit2)
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
