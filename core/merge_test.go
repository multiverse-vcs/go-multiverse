package core

import (
	"context"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestMergeConflicts(t *testing.T) {
	store, err := storage.NewStore(afero.NewMemMapFs(), "/")
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	base, err := Commit(context.TODO(), store, "base")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello\nfoo\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	local, err := Commit(context.TODO(), store, "local", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello\nbar\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	remote, err := Commit(context.TODO(), store, "remote", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	_, err = Merge(context.TODO(), store, base, local, remote)
	if err != nil {
		t.Fatalf("failed to merge %s", err)
	}
}
