package merge

import (
	"context"
	"errors"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// Tree combines the changes to trees a and b onto the base o.
func Tree(ctx context.Context, ds ipld.DAGService, o, a, b cid.Cid) (ipld.Node, error) {
	// fast forward b
	if o == a {
		return object.GetCommitTree(ctx, ds, b)
	}

	// fast forward a
	if o == b {
		return object.GetCommitTree(ctx, ds, a)
	}

	treeO, err := object.GetCommitTree(ctx, ds, o)
	if err != nil {
		return nil, err
	}

	treeA, err := object.GetCommitTree(ctx, ds, a)
	if err != nil {
		return nil, err
	}

	treeB, err := object.GetCommitTree(ctx, ds, b)
	if err != nil {
		return nil, err
	}

	changesA, err := dagutils.Diff(ctx, ds, treeO, treeA)
	if err != nil {
		return nil, err
	}

	changesB, err := dagutils.Diff(ctx, ds, treeO, treeB)
	if err != nil {
		return nil, err
	}

	changes, conflicts := dagutils.MergeDiffs(changesA, changesB)
	for _, c := range conflicts {
		change, err := resolve(ctx, ds, c)
		if err != nil {
			return nil, err
		}

		changes = append(changes, change)
	}

	proto, ok := treeO.(*merkledag.ProtoNode)
	if !ok {
		return nil, errors.New("invalid tree")
	}

	return dagutils.ApplyChange(ctx, ds, proto, changes)
}

// resolve merges the contents of two conflicting dag changes.
func resolve(ctx context.Context, ds ipld.DAGService, c dagutils.Conflict) (*dagutils.Change, error) {
	if c.A.Type == dagutils.Remove {
		return c.B, nil
	}

	if c.B.Type == dagutils.Remove {
		return c.A, nil
	}

	merge, err := File(ctx, ds, c.A.Before, c.A.After, c.B.After)
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
