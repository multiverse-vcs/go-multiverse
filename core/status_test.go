package core

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func (c *Context) TestStatus(t *testing.T) {
	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}

	readme := filepath.Join(c.root, "README.md")
	if err := ioutil.WriteFile(readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write readme file")
	}

	changes, err := c.Status()
	if err != nil {
		t.Fatalf("failed to get status")
	}

	if len(changes) != 1 {
		t.Fatalf("unexpected changes")
	}

	if changes[0].Path != "README.md" {
		t.Fatalf("unexpected change path")
	}

	if changes[0].Type != dagutils.Add {
		t.Fatalf("unexpected change type")
	}
}
