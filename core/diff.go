package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Diff returns a list of changes between the two commit trees.
func Diff(ctx context.Context, dag ipld.DAGService, a, b cid.Cid) ([]*dagutils.Change, error) {
	nodeA, err := dag.Get(ctx, a)
	if err != nil {
		return nil, err
	}

	commitA, err := data.CommitFromCBOR(nodeA.RawData())
	if err != nil {
		return nil, err
	}

	treeA, err := dag.Get(ctx, commitA.Tree)
	if err != nil {
		return nil, err
	}

	nodeB, err := dag.Get(ctx, b)
	if err != nil {
		return nil, err
	}

	commitB, err := data.CommitFromCBOR(nodeB.RawData())
	if err != nil {
		return nil, err
	}

	treeB, err := dag.Get(ctx, commitB.Tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, dag, treeA, treeB)
}
