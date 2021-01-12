package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/spf13/afero"
)

func TestMergeBase(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tree, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	base := data.NewCommit(tree.Cid(), "base")
	baseId, err := data.AddCommit(ctx, dag, base)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	local := data.NewCommit(tree.Cid(), "local", baseId)
	localId, err := data.AddCommit(ctx, dag, local)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	remote := data.NewCommit(tree.Cid(), "remote", baseId)
	remoteId, err := data.AddCommit(ctx, dag, remote)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	mergeId, err := MergeBase(ctx, dag, localId, remoteId)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if mergeId != baseId {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseRemoteAhead(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tree, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	base := data.NewCommit(tree.Cid(), "base")
	baseId, err := data.AddCommit(ctx, dag, base)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	local := data.NewCommit(tree.Cid(), "local", baseId)
	localId, err := data.AddCommit(ctx, dag, local)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	remote := data.NewCommit(tree.Cid(), "remote", localId)
	remoteId, err := data.AddCommit(ctx, dag, remote)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	mergeId, err := MergeBase(ctx, dag, localId, remoteId)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if mergeId != localId {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseLocalAhead(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tree, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	base := data.NewCommit(tree.Cid(), "base")
	baseId, err := data.AddCommit(ctx, dag, base)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	remote := data.NewCommit(tree.Cid(), "remote", baseId)
	remoteId, err := data.AddCommit(ctx, dag, remote)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	local := data.NewCommit(tree.Cid(), "local", remoteId)
	localId, err := data.AddCommit(ctx, dag, local)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	mergeId, err := MergeBase(ctx, dag, localId, remoteId)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if mergeId != remoteId {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseUnrelated(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tree, err := Add(ctx, dag, "", nil)
	if err != nil {
		t.Fatalf("failed to add tree")
	}

	local := data.NewCommit(tree.Cid(), "local")
	localId, err := data.AddCommit(ctx, dag, local)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	remote := data.NewCommit(tree.Cid(), "remote")
	remoteId, err := data.AddCommit(ctx, dag, remote)
	if err != nil {
		t.Fatalf("failed to add commit")
	}

	mergeId, err := MergeBase(ctx, dag, localId, remoteId)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if mergeId.Defined() {
		t.Errorf("uexpected merge base")
	}
}
