package dag

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// Equal returns true if the tree is equal to the tree of the commit with the given id.
func Equal(ctx context.Context, dag ipld.DAGService, id cid.Cid, tree ipld.Node) (bool, error) {
	if !id.Defined() {
		return false, nil
	}

	commit, err := object.GetCommit(ctx, dag, id)
	if err != nil {
		return false, err
	}

	if commit.Tree == tree.Cid() {
		return true, nil
	}

	return false, nil
}
