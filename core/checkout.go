package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Checkout writes the tree of the commit to the root.
func Checkout(ctx context.Context, dag ipld.DAGService, path string, id cid.Cid) error {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return err
	}

	commit, err := data.CommitFromCBOR(node.RawData())
	if err != nil {
		return err
	}

	tree, err := dag.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	return Write(ctx, dag, path, tree)
}
