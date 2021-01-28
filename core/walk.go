package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Walk performs a depth first traversal of parent commits starting at the given id.
func Walk(ctx context.Context, dag ipld.DAGService, id cid.Cid, visit func(cid.Cid) bool) error {
	getLinks := func(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
		commit, err := data.GetCommit(ctx, dag, id)
		if err != nil {
			return nil, err
		}

		return commit.ParentLinks(), nil
	}

	seen := make(map[string]bool)
	wrap := func(id cid.Cid) bool {
		if seen[id.KeyString()] {
			return false
		}

		seen[id.KeyString()] = true
		return visit(id)
	}

	return merkledag.Walk(ctx, getLinks, id, wrap)
}
