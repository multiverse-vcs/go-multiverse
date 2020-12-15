package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestMergeConflicts(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	base, err := Commit(context.TODO(), fs, dag, "base")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello\nfoo\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	local, err := Commit(context.TODO(), fs, dag, "local", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello\nbar\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	remote, err := Commit(context.TODO(), fs, dag, "remote", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	_, err = Merge(context.TODO(), fs, dag, base, local, remote)
	if err != nil {
		t.Fatalf("failed to merge %s", err)
	}
}
