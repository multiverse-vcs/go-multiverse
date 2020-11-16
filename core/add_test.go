package core

import (
	"testing"

	"github.com/ipfs/go-ipfs-files"
)

func TestAddFile(t *testing.T) {
	file := files.NewBytesFile([]byte("foo bar"))

	if _, err := NewMockContext().Add(file); err != nil {
		t.Fatalf("failed to add file: %s", err)
	}
}

func TestAddDir(t *testing.T) {
	file := files.NewMapDirectory(map[string]files.Node{
		"1": files.NewBytesFile([]byte("foo bar")),
		"2": files.NewLinkFile("foo bar", nil),
		"3": files.NewMapDirectory(map[string]files.Node{}),
	})

	if _, err := NewMockContext().Add(file); err != nil {
		t.Fatalf("failed to add dir: %s", err)
	}
}

func TestAddSymlink(t *testing.T) {
	file := files.NewLinkFile("foo", nil)

	if _, err := NewMockContext().Add(file); err != nil {
		t.Fatalf("failed to add symlink: %s", err)
	}
}
