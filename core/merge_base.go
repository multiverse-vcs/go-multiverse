package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// MergeBase returns the best common ancestor of local and remote.
func MergeBase(ctx context.Context, dag ipld.DAGService, local, remote cid.Cid) (cid.Cid, error) {
	history, err := Walk(ctx, dag, local, nil)
	if err != nil {
		return cid.Cid{}, err
	}

	// local is ahead of remote
	if _, ok := history[remote.KeyString()]; ok {
		return remote, nil
	}

	var best cid.Cid
	var err0 error

	// find the least common ancestor by searching
	// for commits that are in both local and remote
	// and that are also independent from each other
	cb := func(id cid.Cid, commit *data.Commit) bool {
		if err0 != nil {
			return false
		}

		if _, ok := history[id.KeyString()]; !ok {
			return true
		}

		var match bool
		if match, err0 = IsAncestor(ctx, dag, best, id); !match {
			best = id
		}

		return false
	}

	if _, err := Walk(ctx, dag, remote, cb); err != nil {
		return cid.Cid{}, err
	}

	return best, err0
}
