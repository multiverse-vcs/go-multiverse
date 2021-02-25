package dag

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// Status returns a list of changes between the tree and commit with the given id.
func Status(ctx context.Context, dag ipld.DAGService, id cid.Cid, tree ipld.Node) (map[string]ChangeType, error) {
	if !id.Defined() {
		return Diff(ctx, dag, &merkledag.ProtoNode{}, tree)
	}

	commit, err := object.GetCommit(ctx, dag, id)
	if err != nil {
		return nil, err
	}

	index, err := dag.Get(ctx, commit.Tree)
	if err != nil {
		return nil, err
	}

	return Diff(ctx, dag, index, tree)
}
