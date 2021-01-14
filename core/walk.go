package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// WalkFun is called for each commit visited by walk.
type WalkFun func(cid.Cid) bool

// Walk performs a depth first traversal of parent commits starting at the given id.
func Walk(ctx context.Context, dag ipld.DAGService, id cid.Cid, cb WalkFun) (map[string]bool, error) {
	getLinks := func(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
		commit, err := data.GetCommit(ctx, dag, id)
		if err != nil {
			return nil, err
		}

		return commit.ParentLinks(), nil
	}

	history := make(map[string]bool)
	visit := func(id cid.Cid) bool {
		if history[id.KeyString()] {
			return false
		}

		history[id.KeyString()] = true
		if cb != nil {
			return cb(id)
		}

		return true
	}

	return history, merkledag.Walk(ctx, getLinks, id, visit)
}
