package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
	"github.com/multiverse-vcs/go-multiverse/object"
)

func TestCommit(t *testing.T) {
	mock := NewMockContext()

	readme := mock.fs.Join(mock.config.Root, "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	idA, err := mock.Commit("first")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	nodeA, err := mock.dag.Get(mock.ctx, idA)
	if err != nil {
		t.Fatalf("failed to get commit")
	}

	commitA, err := object.CommitFromCBOR(nodeA.RawData())
	if err != nil {
		t.Fatalf("failed to decode commit")
	}

	if commitA.Message != "first" {
		t.Errorf("commit message does not match")
	}

	if len(commitA.Parents) != 0 {
		t.Fatalf("commit parent does not match")
	}

	if mock.config.Head != idA {
		t.Errorf("config head does not match")
	}

	idB, err := mock.Commit("second")
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	nodeB, err := mock.dag.Get(mock.ctx, idB)
	if err != nil {
		t.Fatalf("failed to get commit")
	}

	commitB, err := object.CommitFromCBOR(nodeB.RawData())
	if err != nil {
		t.Fatalf("failed to decode commit")
	}

	if commitB.Message != "second" {
		t.Errorf("commit message does not match")
	}

	if len(commitB.Parents) != 1 {
		t.Fatalf("commit parent does not match")
	}

	if commitB.Parents[0] != idA {
		t.Errorf("commit parent does not match")
	}

	if mock.config.Head != idB {
		t.Errorf("config head does not match")
	}
}
