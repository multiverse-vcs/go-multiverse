package core

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// Diff returns a list of changes between the two commit trees.
func Diff(ctx context.Context, store *storage.Store, a, b cid.Cid) ([]*dagutils.Change, error) {
	nodeA, err := store.Dag.Get(ctx, a)
	if err != nil {
		return nil, err
	}

	commitA, err := object.CommitFromCBOR(nodeA.RawData())
	if err != nil {
		return nil, err
	}

	treeA, err := store.Dag.Get(ctx, commitA.Tree)
	if err != nil {
		return nil, err
	}

	nodeB, err := store.Dag.Get(ctx, b)
	if err != nil {
		return nil, err
	}

	commitB, err := object.CommitFromCBOR(nodeB.RawData())
	if err != nil {
		return nil, err
	}

	treeB, err := store.Dag.Get(ctx, commitB.Tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, store.Dag, treeA, treeB)
}
