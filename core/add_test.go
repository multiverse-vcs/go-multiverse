package core

import (
	"testing"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-unixfs/file"
)

func TestAdd(t *testing.T) {
	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}

	test := files.NewMapDirectory(map[string]files.Node{
		"1": files.NewBytesFile([]byte("foo bar")),
		"2": files.NewLinkFile("foo bar", nil),
		"3": files.NewMapDirectory(map[string]files.Node{}),
	})

	node, err := c.Add(test)
	if err != nil {
		t.Fatalf("failed to add file: %s", err)
	}

	file, err := unixfile.NewUnixfsFile(c.ctx, c.dag, node)
	if err != nil {
		t.Fatalf("failed to read unixfile: %s", err)
	}

	dir, ok := file.(files.Directory)
	if !ok {
		t.Fatalf("expected file to be a directory")
	}

	entries := dir.Entries()
	if !entries.Next() {
		t.Fatalf("unexpected entries")
	}

	if entries.Name() != "1" {
		t.Errorf("unexpected entry")
	}

	if !entries.Next() {
		t.Fatalf("unexpected entries")
	}

	if entries.Name() != "2" {
		t.Errorf("unexpected entry")
	}

	if !entries.Next() {
		t.Fatalf("unexpected entries")
	}

	if entries.Name() != "3" {
		t.Errorf("unexpected entry")
	}

	if entries.Next() {
		t.Errorf("unexpected entries")
	}
}
