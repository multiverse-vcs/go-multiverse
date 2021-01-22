package core

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

// Merge combines the work trees of a and b into the base o.
func Merge(ctx context.Context, dag ipld.DAGService, o, a, b cid.Cid) (ipld.Node, error) {
	changesA, err := Diff(ctx, dag, o, a)
	if err != nil {
		return nil, err
	}

	changesB, err := Diff(ctx, dag, o, b)
	if err != nil {
		return nil, err
	}

	changes, conflicts := dagutils.MergeDiffs(changesA, changesB)
	for _, c := range conflicts {
		change, err := conflict(ctx, dag, c)
		if err != nil {
			return nil, err
		}

		changes = append(changes, change)
	}

	base, err := data.GetCommit(ctx, dag, o)
	if err != nil {
		return nil, err
	}

	tree, err := dag.Get(ctx, base.Tree)
	if err != nil {
		return nil, err
	}

	proto, ok := tree.(*merkledag.ProtoNode)
	if !ok {
		return nil, errors.New("invalid commit tree")
	}

	return dagutils.ApplyChange(ctx, dag, proto, changes)
}

// conflict combines the contents of two conflicting dag changes.
func conflict(ctx context.Context, dag ipld.DAGService, c dagutils.Conflict) (*dagutils.Change, error) {
	if c.A.Type == dagutils.Remove {
		return c.B, nil
	}

	if c.B.Type == dagutils.Remove {
		return c.A, nil
	}

	merge, err := unixfs.Merge(ctx, dag, c.A.Before, c.A.After, c.B.After)
	if err != nil {
		return nil, err
	}

	change := dagutils.Mod
	if c.A.Type == dagutils.Add && c.B.Type == dagutils.Add {
		change = dagutils.Add
	}

	return &dagutils.Change{
		Type:   change,
		Path:   c.A.Path,
		Before: c.A.Before,
		After:  merge.Cid(),
	}, nil
}
