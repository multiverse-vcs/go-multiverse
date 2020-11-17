package core

import (
	"testing"

	"github.com/ipfs/go-cid"
)

func TestCommit(t *testing.T) {
	parent, err := cid.Parse("bagaybqabciqchvwlfbygi7w76xf5q64so7lqavtwepcdt4bp4t7ha42ldxa2sya")
	if err != nil {
		t.Fatalf("failed to parse parent cid")
	}

	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}
	c.config.Head = parent

	commit, err := c.Commit("foo bar")
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	if commit.Message != "foo bar" {
		t.Errorf("commit message does not match")
	}

	if len(commit.Parents) != 1 {
		t.Fatalf("commit parent does not match")
	}

	if commit.Parents[0] != parent {
		t.Errorf("commit parent does not match")
	}

	if c.config.Head != commit.Cid() {
		t.Errorf("config head does not match")
	}
}
