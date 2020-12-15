package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestCheckout(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	id, err := Commit(context.TODO(), fs, dag, "init")
	if err != nil {
		t.Fatalf("failed to create worktree")
	}

	if err := fs.RemoveAll(""); err != nil {
		t.Fatalf("failed to remove all")
	}

	if err := Checkout(context.TODO(), fs, dag, id); err != nil {
		t.Fatalf("failed to checkout")
	}

	if _, err := fs.Stat("README.md"); err != nil {
		t.Fatalf("failed to lstat file")
	}
}
