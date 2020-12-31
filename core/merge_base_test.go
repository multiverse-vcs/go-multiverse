package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestMergeBase(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	base, err := Commit(ctx, dag, "/", "base")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(ctx, dag, "/", "local", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(ctx, dag, "/", "remote", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(ctx, dag, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != base {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseRemoteAhead(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	base, err := Commit(ctx, dag, "/", "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(ctx, dag, "/", "local", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(ctx, dag, "/", "remote", local)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(ctx, dag, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != local {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseLocalAhead(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	base, err := Commit(ctx, dag, "/", "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(ctx, dag, "/", "remote", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(ctx, dag, "/", "local", remote)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(ctx, dag, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != remote {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseUnrelated(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	local, err := Commit(ctx, dag, "/", "local")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(ctx, dag, "/", "remote")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(ctx, dag, local, remote)
	if merge.Defined() {
		t.Errorf("uexpected merge base")
	}
}
