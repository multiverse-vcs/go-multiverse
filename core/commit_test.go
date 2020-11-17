package core

import (
	"testing"
)

func TestCommit(t *testing.T) {
	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}

	first, err := c.Commit("first")
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	if first.Message != "first" {
		t.Errorf("commit message does not match")
	}

	if len(first.Parents) != 0 {
		t.Fatalf("commit parent does not match")
	}

	if c.config.Head != first.Cid() {
		t.Errorf("config head does not match")
	}

	second, err := c.Commit("second")
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	if second.Message != "second" {
		t.Errorf("commit message does not match")
	}

	if len(second.Parents) != 1 {
		t.Fatalf("commit parent does not match")
	}

	if second.Parents[0] != first.Cid() {
		t.Errorf("commit parent does not match")
	}

	if c.config.Head != second.Cid() {
		t.Errorf("config head does not match")
	}
}
