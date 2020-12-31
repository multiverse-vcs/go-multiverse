package core

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Status returns a list of changes between the worktree and commit with the given id.
func Status(ctx context.Context, dag ipld.DAGService, path string, id cid.Cid) (map[string]dagutils.ChangeType, error) {
	tree, err := Worktree(ctx, dag, path)
	if err != nil {
		return nil, err
	}

	if !id.Defined() {
		return mapChanges(ctx, dag, &merkledag.ProtoNode{}, tree)
	}

	node, err := dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	commit, err := data.CommitFromCBOR(node.RawData())
	if err != nil {
		return nil, err
	}

	nodeA, err := dag.Get(ctx, commit.Tree)
	if err != nil {
		return nil, err
	}

	return mapChanges(ctx, dag, nodeA, tree)
}

// mapChanges returns a map of unique file changes.
func mapChanges(ctx context.Context, dag ipld.DAGService, nodeA, nodeB ipld.Node) (map[string]dagutils.ChangeType, error) {
	changes, err := dagutils.Diff(ctx, dag, nodeA, nodeB)
	if err != nil {
		return nil, err
	}

	diffs := make(map[string]dagutils.ChangeType)
	for _, change := range changes {
		if change.Path == "" {
			continue
		}

		if _, ok := diffs[change.Path]; ok {
			diffs[change.Path] = dagutils.Mod
		} else {
			diffs[change.Path] = change.Type
		}
	}

	return diffs, nil
}
