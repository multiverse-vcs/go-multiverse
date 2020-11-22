package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
	ufsio "github.com/ipfs/go-unixfs/io"
)

func TestWorktree(t *testing.T) {
	mock := NewMockContext()

	IgnoreRules = []string{"*.exe"}

	exe := mock.Fs.Join(mock.Fs.Root(), "test.exe")
	if err := fsutil.WriteFile(mock.Fs, exe, []byte{0, 0, 0}, 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	readme := mock.Fs.Join(mock.Fs.Root(), "README.md")
	if err := fsutil.WriteFile(mock.Fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := mock.Worktree()
	if err != nil {
		t.Fatalf("failed to create worktree")
	}

	dir, err := ufsio.NewDirectoryFromNode(mock.Dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = dir.Find(mock, "README.md")
	if err != nil {
		t.Errorf("failed to find file")
	}

	_, err = dir.Find(mock, "test.exe")
	if err == nil {
		t.Errorf("expected file to be ignored")
	}
}
