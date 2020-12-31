package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestMergeConflicts(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	base, err := Commit(ctx, dag, "", "base")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello\nfoo\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	local, err := Commit(ctx, dag, "", "local", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(fs, "README.md", []byte("hello\nbar\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	remote, err := Commit(ctx, dag, "", "remote", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	_, err = Merge(ctx, dag, base, local, remote)
	if err != nil {
		t.Fatalf("failed to merge %s", err)
	}
}
