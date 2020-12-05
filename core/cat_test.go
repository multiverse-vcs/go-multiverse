package core

import (
	"context"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestCat(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, "test.txt", []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), store, "test.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	text, err := Cat(context.TODO(), store, node.Cid())
	if err != nil {
		t.Fatal("failed to cat file")
	}

	if text != "foo bar" {
		t.Error("unexpected file contents")
	}
}
