package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/spf13/afero"
)

func TestWorktree(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	IgnoreRules = []string{"*.exe"}
	if err := afero.WriteFile(fs, "test.exe", []byte{0, 0, 0}, 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Worktree(context.TODO(), fs, dag)
	if err != nil {
		t.Fatalf("failed to create worktree")
	}

	dir, err := ufsio.NewDirectoryFromNode(dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = dir.Find(context.TODO(), "README.md")
	if err != nil {
		t.Errorf("failed to find file")
	}

	_, err = dir.Find(context.TODO(), "test.exe")
	if err == nil {
		t.Errorf("expected file to be ignored")
	}
}
