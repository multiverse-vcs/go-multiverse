package core

import (
	"testing"

	"github.com/multiverse-vcs/go-multiverse/object"
)

func TestCommit(t *testing.T) {
	mock := NewMockContext()

	if err := mock.Fs.MkdirAll(mock.Fs.Root(), 0755); err != nil {
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

	node, err := mock.Dag.Get(mock, id)
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

	if mock.Config.Head != id {
		t.Errorf("config head does not match")
	}

	if mock.Config.Base != id {
		t.Errorf("config base does not match")
	}
}

func TestCommitDetached(t *testing.T) {
	mock := NewMockContext()

	if err := mock.Fs.MkdirAll(mock.Fs.Root(), 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	commit, err := mock.Commit("init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	if _, err := mock.Commit("second"); err != nil {
		t.Fatalf("failed to create commit")
	}

	mock.Config.Base = commit
	if _, err := mock.Commit("detached"); err == nil {
		t.Errorf("expected commit error")
	}
}
