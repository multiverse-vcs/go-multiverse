package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
)

// MergeBase returns the best common ancestor of local and remote.
func MergeBase(ctx context.Context, dag ipld.DAGService, local, remote cid.Cid) (cid.Cid, error) {
	history := make(map[string]bool)
	visit := func(id cid.Cid) bool {
		history[id.KeyString()] = true
		return true
	}

	if err := Walk(ctx, dag, local, visit); err != nil {
		return cid.Cid{}, err
	}

	// local is ahead of remote
	if history[remote.KeyString()] {
		return remote, nil
	}

	var best cid.Cid
	var err0 error
	var match bool

	// find the least common ancestor by searching
	// for commits that are in both local and remote
	// and that are also independent from each other
	visit = func(id cid.Cid) bool {
		if err0 != nil {
			return false
		}

		if !history[id.KeyString()] {
			return true
		}

		if match, err0 = IsAncestor(ctx, dag, best, id); !match {
			best = id
		}

		return false
	}

	if err := Walk(ctx, dag, remote, visit); err != nil {
		return cid.Cid{}, err
	}

	return best, err0
}

// IsAncestor returns true if child is an ancestor of parent.
func IsAncestor(ctx context.Context, dag ipld.DAGService, parent, child cid.Cid) (bool, error) {
	if !parent.Defined() || !child.Defined() {
		return false, nil
	}

	var match bool
	visit := func(id cid.Cid) bool {
		match = (id == child)
		return !match
	}

	if err := Walk(ctx, dag, parent, visit); err != nil {
		return false, err
	}

	return match, nil
}
