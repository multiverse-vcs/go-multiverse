package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// IsAncestor returns true if child is an ancestor of parent.
func IsAncestor(ctx context.Context, dag ipld.DAGService, parent, child cid.Cid) (bool, error) {
	if !(parent.Defined() && child.Defined()) {
		return false, nil
	}

	cb := func(id cid.Cid, commit *data.Commit) bool {
		return id != child
	}

	history, err := Walk(ctx, dag, parent, cb)
	if err != nil {
		return false, err
	}

	_, ok := history[child.KeyString()]
	return ok, nil
}
