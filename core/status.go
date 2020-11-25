package core

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// Status returns a list of changes between the worktree and commit with the given id.
func Status(ctx context.Context, store *storage.Store, id cid.Cid) ([]*dagutils.Change, error) {
	tree, err := Worktree(ctx, store)
	if err != nil {
		return nil, err
	}

	if !id.Defined() {
		return dagutils.Diff(ctx, store.Dag, &merkledag.ProtoNode{}, tree)
	}

	node, err := store.Dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return nil, err
	}

	nodeA, err := store.Dag.Get(ctx, commit.Tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, store.Dag, nodeA, tree)
}
