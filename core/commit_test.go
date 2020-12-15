package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/spf13/afero"
)

func TestCommit(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	parent, err := Commit(context.TODO(), fs, dag, "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	id, err := Commit(context.TODO(), fs, dag, "changes", parent)
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	node, err := dag.Get(context.TODO(), id)
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
