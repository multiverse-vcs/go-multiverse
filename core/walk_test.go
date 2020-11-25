package core

import (
	"context"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestWalk(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	idA, err := Commit(context.TODO(), store, "first")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	idB, err := Commit(context.TODO(), store, "second", idA)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	history, err := Walk(context.TODO(), store, idB, nil)
	if err != nil {
		t.Fatalf("failed to walk")
	}

	if len(history) != 2 {
		t.Fatalf("cids do not match")
	}

	if _, ok := history[idA.KeyString()]; !ok {
		t.Errorf("cids do not match")
	}

	if _, ok := history[idB.KeyString()]; !ok {
		t.Errorf("cids do not match")
	}
}
