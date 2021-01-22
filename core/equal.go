package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

// Equal returns true if the worktree is equal to the tree of the commit.
func Equal(ctx context.Context, dag ipld.DAGService, path string, ignore unixfs.Ignore, id cid.Cid) (bool, error) {
	mem := dagutils.NewMemoryDagService()

	tree, err := unixfs.Add(ctx, mem, path, ignore)
	if err != nil {
		return false, err
	}

	head, err := data.GetCommit(ctx, dag, id)
	if err != nil {
		return false, err
	}

	if tree.Cid() != head.Tree {
		return false, nil
	}

	return true, nil
}
