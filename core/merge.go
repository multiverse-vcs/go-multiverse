package core

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/go-ipld-multiverse"
)

// Merge combines the repo histories of the local and remote commits.
func (c *Core) Merge(ctx context.Context, ref path.Path) error {
	p, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return err
	}

	if p.Cid().Type() != ipldmulti.CommitCodec {
		return ErrInvalidRef
	}

	base, err := c.MergeBase(ctx, c.Config.Head, p.Cid())
	if err != nil {
		return nil
	}

	fmt.Println(base.Cid())
	return nil
}

// MergeBase returns the best common ancestor for merging.
func (c *Core) MergeBase(ctx context.Context, local, remote cid.Cid) (*ipldmulti.Commit, error) {
	history, err := c.History(ctx, local)
	if err != nil {
		return nil, err
	}

	if history[remote.KeyString()] {
		return nil, ErrMergeAhead
	}

	var filter CommitFilter = func(commit *ipldmulti.Commit) bool {
		return history[commit.Cid().KeyString()]
	}

	iter := c.NewCommitIter(remote).WithFilter(&filter).WithLimit(&filter)

	base, err := iter.Next(ctx)
	if err != nil {
		return nil, err
	}

	if base == nil {
		return nil, ErrMergeBase
	}

	return base, iter.ForEach(ctx, func(commit *ipldmulti.Commit) error {
		ancestor, err := c.IsAncestor(ctx, commit.Cid(), base.Cid())
		if err != nil {
			return err
		}

		if !ancestor {
			base = commit
		}

		return nil
	})
}

// IsAncestor checks if child is an ancestor of parent.
func (c *Core) IsAncestor(ctx context.Context, child, parent cid.Cid) (bool, error) {
	var filter CommitFilter = func(commit *ipldmulti.Commit) bool {
		return commit.Cid().Equals(child)
	}

	commit, err := c.NewCommitIter(parent).WithFilter(&filter).Next(ctx)
	if err != nil {
		return false, err
	}

	if commit == nil {
		return false, nil
	}

	return true, nil
}

// History returns a flattened map of the commit history.
func (c *Core) History(ctx context.Context, id cid.Cid) (map[string]bool, error) {
	history := make(map[string]bool)
	iter := c.NewCommitIter(c.Config.Head)

	return history, iter.ForEach(ctx, func(commit *ipldmulti.Commit) error {
		history[commit.Cid().KeyString()] = true
		return nil
	})
}

