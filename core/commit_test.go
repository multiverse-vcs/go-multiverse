package core

import (
	"context"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

func TestCommit(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	parent, err := Commit(context.TODO(), store, "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	id, err := Commit(context.TODO(), store, "changes", parent)
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	node, err := store.Dag.Get(context.TODO(), id)
	if err != nil {
		t.Fatalf("failed to get commit")
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		t.Fatalf("failed to decode commit")
	}

	if commit.Message != "changes" {
		t.Errorf("commit message does not match")
	}

	if len(commit.Parents) != 1 {
		t.Fatalf("commit parent does not match")
	}

	if commit.Parents[0] != parent {
		t.Errorf("commit parent does not match")
	}
}
