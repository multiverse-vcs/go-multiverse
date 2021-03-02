package dag

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// IsAncestor returns true if child is an ancestor of parent.
func IsAncestor(ctx context.Context, ds ipld.NodeGetter, parent, child cid.Cid) (bool, error) {
	if !parent.Defined() || !child.Defined() {
		return false, nil
	}

	var match bool
	visit := func(id cid.Cid) bool {
		match = (id == child)
		return !match
	}

	if err := Walk(ctx, ds, parent, visit); err != nil {
		return false, err
	}

	return match, nil
}

type walker struct {
	ipld.NodeGetter
}

func (w walker) getLinks(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
	commit, err := object.GetCommit(ctx, w, id)
	if err != nil {
		return nil, err
	}

	return commit.ParentLinks(), nil
}

// Walk performs a depth first traversal of parent commits starting at the given id.
func Walk(ctx context.Context, ds ipld.NodeGetter, id cid.Cid, visit func(cid.Cid) bool) error {
	return merkledag.Walk(ctx, walker{ds}.getLinks, id, visit)
}
