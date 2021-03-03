package merge

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
)

// Base returns the best common ancestor of local and remote.
func Base(ctx context.Context, ds ipld.NodeGetter, local, remote cid.Cid) (cid.Cid, error) {
	if !local.Defined() || !remote.Defined() {
		return cid.Cid{}, nil
	}

	refs := cid.NewSet()
	if err := dag.Walk(ctx, ds, local, refs.Visit); err != nil {
		return cid.Cid{}, err
	}

	if refs.Has(remote) {
		return remote, nil
	}

	var best cid.Cid
	var err0 error
	var match bool

	// find the least common ancestor by searching
	// for commits that are in both local and remote
	// and that are also independent from each other
	visit := func(id cid.Cid) bool {
		if err0 != nil {
			return false
		}

		if !refs.Has(id) {
			return true
		}

		if match, err0 = dag.IsAncestor(ctx, ds, best, id); !match {
			best = id
		}

		return false
	}

	if err := dag.Walk(ctx, ds, remote, visit); err != nil {
		return cid.Cid{}, err
	}

	return best, err0
}
