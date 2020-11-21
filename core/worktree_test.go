package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
	ufsio "github.com/ipfs/go-unixfs/io"
)

func TestWorktree(t *testing.T) {
	mock := NewMockContext()

	readme := mock.fs.Join(mock.fs.Root(), "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	dot := mock.fs.Join(mock.fs.Root(), ".multiverse")
	if err := mock.fs.MkdirAll(dot, 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	node, err := mock.Worktree()
	if err != nil {
		t.Fatalf("failed to create worktree")
	}

	dir, err := ufsio.NewDirectoryFromNode(mock.dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = dir.Find(mock.ctx, "README.md")
	if err != nil {
		t.Errorf("failed to find file")
	}

	_, err = dir.Find(mock.ctx, ".multiverse")
	if err == nil {
		t.Errorf("expected file to be ignored")
	}
}

func TestWorktreeEmpty(t *testing.T) {
	mock := NewMockContext()

	if err := mock.fs.MkdirAll(mock.fs.Root(), 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	if _, err := mock.Worktree(); err != nil {
		t.Fatalf("failed to create worktree")
	}
}
