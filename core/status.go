package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Status returns a list of changes between the worktree and commit with the given id.
func Status(ctx context.Context, dag ipld.DAGService, path string, filter Filter, id cid.Cid) ([]*dagutils.Change, error) {
	mem := dagutils.NewMemoryDagService()

	tree, err := Add(ctx, mem, path, filter)
	if err != nil {
		return nil, err
	}

	if !id.Defined() {
		return dagutils.Diff(ctx, mem, &merkledag.ProtoNode{}, tree)
	}

	node, err := dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	commit, err := data.CommitFromCBOR(node.RawData())
	if err != nil {
		return nil, err
	}

	nodeA, err := dag.Get(ctx, commit.Tree)
	if err != nil {
		return nil, err
	}

	var ids []cid.Cid
	visit := func(id cid.Cid) bool {
		ids = append(ids, id)
		return true
	}

	getLinks := merkledag.GetLinksWithDAG(dag)
	if err := merkledag.Walk(ctx, getLinks, commit.Tree, visit); err != nil {
		return nil, err
	}

	for opt := range dag.GetMany(ctx, ids) {
		if opt.Err != nil {
			return nil, opt.Err
		}

		if err := mem.Add(ctx, opt.Node); err != nil {
			return nil, err
		}
	}

	return dagutils.Diff(ctx, mem, nodeA, tree)
}
