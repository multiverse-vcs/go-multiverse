package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/spf13/afero"
)

func TestCommit(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	parent, err := Commit(ctx, dag, "/", "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	id, err := Commit(ctx, dag, "/", "changes", parent)
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	node, err := dag.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get commit")
	}

	commit, err := data.CommitFromCBOR(node.RawData())
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
