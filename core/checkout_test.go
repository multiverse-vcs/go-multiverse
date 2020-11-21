package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
)

func TestCheckout(t *testing.T) {
	mock := NewMockContext()

	readme := mock.fs.Join(mock.fs.Root(), "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	id, err := mock.Commit("init")
	if err != nil {
		t.Fatalf("failed to create worktree")
	}

	if err := fsutil.RemoveAll(mock.fs, mock.fs.Root()); err != nil {
		t.Fatalf("failed to remove all")
	}

	if err := mock.Checkout(id); err != nil {
		t.Fatalf("failed to checkout")
	}

	if _, err := mock.fs.Lstat(readme); err != nil {
		t.Fatalf("failed to lstat file")
	}
}
