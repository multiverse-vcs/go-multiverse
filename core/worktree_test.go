package core

import (
	"context"
	"testing"

	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestWorktree(t *testing.T) {
	store, err := storage.NewStore(afero.NewMemMapFs(), "/")
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	IgnoreRules = []string{"*.exe"}
	if err := afero.WriteFile(store.Cwd, "test.exe", []byte{0, 0, 0}, 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Worktree(context.TODO(), store)
	if err != nil {
		t.Fatalf("failed to create worktree")
	}

	dir, err := ufsio.NewDirectoryFromNode(store.Dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = dir.Find(context.TODO(), "README.md")
	if err != nil {
		t.Errorf("failed to find file")
	}

	_, err = dir.Find(context.TODO(), "test.exe")
	if err == nil {
		t.Errorf("expected file to be ignored")
	}
}
