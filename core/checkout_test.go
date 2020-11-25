package core

import (
	"context"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestCheckout(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	id, err := Commit(context.TODO(), store, "init")
	if err != nil {
		t.Fatalf("failed to create worktree")
	}

	if err := store.Cwd.RemoveAll(""); err != nil {
		t.Fatalf("failed to remove all")
	}

	if err := Checkout(context.TODO(), store, id); err != nil {
		t.Fatalf("failed to checkout")
	}

	if _, err := store.Cwd.Stat("README.md"); err != nil {
		t.Fatalf("failed to lstat file")
	}
}
