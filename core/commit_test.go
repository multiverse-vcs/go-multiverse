package core

import (
	"testing"

	//fsutil "github.com/go-git/go-billy/v5/util"
	"github.com/multiverse-vcs/go-multiverse/object"
)

func TestCommit(t *testing.T) {
	mock := NewMockContext()

	if err := mock.fs.MkdirAll(mock.fs.Root(), 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	parent, err := mock.Commit("init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	id, err := mock.Commit("changes")
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	node, err := mock.dag.Get(mock.ctx, id)
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

	if mock.cfg.Head != id {
		t.Errorf("config head does not match")
	}

	if mock.cfg.Base != id {
		t.Errorf("config base does not match")
	}
}
