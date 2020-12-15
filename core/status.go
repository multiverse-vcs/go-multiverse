package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/spf13/afero"
)

// Status returns a list of changes between the worktree and commit with the given id.
func Status(ctx context.Context, fs afero.Fs, dag ipld.DAGService, id cid.Cid) ([]*dagutils.Change, error) {
	tree, err := Worktree(ctx, fs, dag)
	if err != nil {
		return nil, err
	}

	if !id.Defined() {
		return dagutils.Diff(ctx, dag, &merkledag.ProtoNode{}, tree)
	}

	node, err := dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return nil, err
	}

	nodeA, err := dag.Get(ctx, commit.Tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, dag, nodeA, tree)
}
