package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
)

func TestWalk(t *testing.T) {
	mock := NewMockContext()

	readme := mock.fs.Join(mock.config.Root, "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	idA, err := mock.Commit("first")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	idB, err := mock.Commit("second")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	history, err := mock.Walk(idB, nil)
	if err != nil {
		t.Fatalf("failed to walk")
	}

	if len(history) != 2 {
		t.Fatalf("cids do not match")
	}

	if _, ok := history[idA.KeyString()]; !ok {
		t.Errorf("cids do not match")
	}

	if _, ok := history[idB.KeyString()]; !ok {
		t.Errorf("cids do not match")
	}
}
