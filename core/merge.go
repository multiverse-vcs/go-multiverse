package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/diff"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// Merge combines the work trees of a and b o
func Merge(ctx context.Context, store *storage.Store, local, remote cid.Cid) error {
	base, err := MergeBase(ctx, store, local, remote)
	if err != nil {
		return err
	}

	if base == local {
		return Checkout(ctx, store, remote)
	}

	if base == remote {
		return Checkout(ctx, store, local)
	}

	changesA, err := Diff(ctx, store, base, local)
	if err != nil {
		return err
	}

	changesB, err := Diff(ctx, store, base, remote)
	if err != nil {
		return err
	}

	changes, conflicts := dagutils.MergeDiffs(changesA, changesB)
	for _, c := range conflicts {
		change, err := mergeConflict(ctx, store, c.A, c.B)
		if err != nil {
			return err
		}

		changes = append(changes, change)
	}

	node, err := store.Dag.Get(ctx, base)
	if err != nil {
		return err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return err
	}

	tree, err := store.Dag.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	proto, ok := tree.(*merkledag.ProtoNode)
	if !ok {
		return errors.New("invalid commit tree")
	}

	merge, err := dagutils.ApplyChange(ctx, store.Dag, proto, changes)
	if err != nil {
		return err
	}

	return Write(ctx, store, "", merge)
}

// mergeConflict combines the contents of two conflicting dag changes.
func mergeConflict(ctx context.Context, store *storage.Store, a, b *dagutils.Change) (*dagutils.Change, error) {
	if a.Type == dagutils.Remove {
		return b, nil
	}

	if b.Type == dagutils.Remove {
		return a, nil
	}

	textO, err := Cat(ctx, store, a.Before)
	if err != nil {
		return nil, err
	}

	textA, err := Cat(ctx, store, a.After)
	if err != nil {
		return nil, err
	}

	textB, err := Cat(ctx, store, b.After)
	if err != nil {
		return nil, err
	}

	merged := diff.Merge(textO, textA, textB)
	reader := strings.NewReader(merged)

	merge, err := add(ctx, store, reader)
	if err != nil {
		return nil, err
	}

	change := dagutils.Change{
		Type:   dagutils.Mod,
		Path:   a.Path,
		Before: a.Before,
		After:  merge.Cid(),
	}

	if a.Type == dagutils.Add && b.Type == dagutils.Add {
		change.Type = dagutils.Add
	}

	return &change, nil
}
