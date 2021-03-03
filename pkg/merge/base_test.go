package merge

import (
	"context"
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

func TestBase(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()

	tree, err := cid.Decode("QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ")
	if err != nil {
		t.Fatal("failed to decode cid")
	}

	base := object.NewCommit()
	base.Tree = tree
	base.Message = "base"

	baseID, err := object.AddCommit(ctx, mem, base)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	local := object.NewCommit()
	local.Tree = tree
	local.Message = "local"
	local.Parents = []cid.Cid{baseID}

	localID, err := object.AddCommit(ctx, mem, local)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	remote := object.NewCommit()
	remote.Tree = tree
	remote.Message = "remote"
	remote.Parents = []cid.Cid{baseID}

	remoteID, err := object.AddCommit(ctx, mem, remote)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	mergeID, err := Base(ctx, mem, localID, remoteID)
	if err != nil {
		t.Fatal("failed to get merge base")
	}

	if mergeID != baseID {
		t.Error("unexpected merge base")
	}
}

func TestBaseRemoteAhead(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()

	tree, err := cid.Decode("QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ")
	if err != nil {
		t.Fatal("failed to decode cid")
	}

	base := object.NewCommit()
	base.Tree = tree
	base.Message = "base"

	baseID, err := object.AddCommit(ctx, mem, base)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	local := object.NewCommit()
	local.Tree = tree
	local.Message = "local"
	local.Parents = []cid.Cid{baseID}

	localID, err := object.AddCommit(ctx, mem, local)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	remote := object.NewCommit()
	remote.Tree = tree
	remote.Message = "remote"
	remote.Parents = []cid.Cid{localID}

	remoteID, err := object.AddCommit(ctx, mem, remote)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	mergeID, err := Base(ctx, mem, localID, remoteID)
	if err != nil {
		t.Fatal("failed to get merge base")
	}

	if mergeID != localID {
		t.Error("unexpected merge base")
	}
}

func TestBaseLocalAhead(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()

	tree, err := cid.Decode("QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ")
	if err != nil {
		t.Fatal("failed to decode cid")
	}

	base := object.NewCommit()
	base.Tree = tree
	base.Message = "base"

	baseID, err := object.AddCommit(ctx, mem, base)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	remote := object.NewCommit()
	remote.Tree = tree
	remote.Message = "remote"
	remote.Parents = []cid.Cid{baseID}

	remoteID, err := object.AddCommit(ctx, mem, remote)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	local := object.NewCommit()
	local.Tree = tree
	local.Message = "local"
	local.Parents = []cid.Cid{remoteID}

	localID, err := object.AddCommit(ctx, mem, local)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	mergeID, err := Base(ctx, mem, localID, remoteID)
	if err != nil {
		t.Fatal("failed to get merge base")
	}

	if mergeID != remoteID {
		t.Error("unexpected merge base")
	}
}

func TestBaseUnrelated(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()

	tree, err := cid.Decode("QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ")
	if err != nil {
		t.Fatal("failed to decode cid")
	}

	local := object.NewCommit()
	local.Tree = tree
	local.Message = "local"

	localID, err := object.AddCommit(ctx, mem, local)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	remote := object.NewCommit()
	remote.Tree = tree
	remote.Message = "remote"

	remoteID, err := object.AddCommit(ctx, mem, remote)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	mergeID, err := Base(ctx, mem, localID, remoteID)
	if err != nil {
		t.Fatal("failed to get merge base")
	}

	if mergeID.Defined() {
		t.Error("uexpected merge base")
	}
}
