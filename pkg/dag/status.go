package dag

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// Status returns the changes between the given tree and commit with id.
func Status(ctx context.Context, ds ipld.DAGService, tree ipld.Node, id cid.Cid) (map[string]dagutils.ChangeType, error) {
	if !id.Defined() {
		return Diff(ctx, ds, &merkledag.ProtoNode{}, tree)
	}

	commit, err := object.GetCommit(ctx, ds, id)
	if err != nil {
		return nil, err
	}

	index, err := ds.Get(ctx, commit.Tree)
	if err != nil {
		return nil, err
	}

	return Diff(ctx, ds, index, tree)
}

// Diff returns a flattened map of changes between before and after.
func Diff(ctx context.Context, ds ipld.DAGService, before, after ipld.Node) (map[string]dagutils.ChangeType, error) {
	changes, err := dagutils.Diff(ctx, ds, before, after)
	if err != nil {
		return nil, err
	}

	diffs := make(map[string]dagutils.ChangeType)
	for _, change := range changes {
		if _, ok := diffs[change.Path]; ok {
			diffs[change.Path] = dagutils.Mod
		} else if change.Path != "" {
			diffs[change.Path] = change.Type
		}
	}

	return diffs, nil
}
