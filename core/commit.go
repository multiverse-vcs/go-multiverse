package core

import (
	"context"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Commit creates a new commit.
func Commit(ctx context.Context, dag ipld.DAGService, path string, message string, parents ...cid.Cid) (cid.Cid, error) {
	tree, err := Worktree(ctx, dag, path)
	if err != nil {
		return cid.Cid{}, err
	}

	commit := data.Commit{
		Date:    time.Now(),
		Message: message,
		Tree:    tree.Cid(),
		Parents: parents,
	}

	node, err := cbornode.WrapObject(&commit, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := dag.Add(ctx, node); err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}
