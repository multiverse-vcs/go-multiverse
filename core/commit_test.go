package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
)

func TestCommit(t *testing.T) {
	mock := NewMockContext()

	readme := mock.fs.Join(mock.config.Root, "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	first, err := mock.Commit("first")
	if err != nil {
		t.Fatalf("failed to commit: %s", err)
	}

	if first.Message != "first" {
		t.Errorf("commit message does not match")
	}

	if len(first.Parents) != 0 {
		t.Fatalf("commit parent does not match")
	}

	if mock.config.Head != first.Cid() {
		t.Errorf("config head does not match")
	}

	second, err := mock.Commit("second")
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

	if mock.config.Head != second.Cid() {
		t.Errorf("config head does not match")
	}
}
