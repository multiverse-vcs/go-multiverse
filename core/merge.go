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

	bases, err := c.MergeBase(ctx, c.Config.Head, p.Cid())
	if err != nil {
		return err
	}

	if len(bases) == 0 {
		return ErrMergeBase
	}

	fmt.Println(bases[0].Cid().String())
	return nil
}

// MergeBase returns a list of possible merge bases for local and remote.
func (c *Core) MergeBase(ctx context.Context, local, remote cid.Cid) ([]*ipldmulti.Commit, error) {
	history, err := c.NewHistory(local).Flatten(ctx)
	if err != nil {
		return nil, err
	}

	if history[remote.KeyString()] {
		return nil, ErrMergeAhead
	}

	var filter HistoryFilter = func(commit *ipldmulti.Commit) bool {
		return history[commit.Cid().KeyString()]
	}

	bases := make([]*ipldmulti.Commit, 0)

	var callback HistoryCallback = func(commit *ipldmulti.Commit) error {
		bases = append(bases, commit)
		return nil
	}

	return bases, c.NewFilterHistory(remote, &filter, &filter).ForEach(ctx, callback)
}

// IsAncestor checks if child is an ancestor of parent.
func (c *Core) IsAncestor(ctx context.Context, child, parent cid.Cid) (bool, error) {
	var filter HistoryFilter = func(commit *ipldmulti.Commit) bool {
		return commit.Cid().Equals(child)
	}

	commit, err := c.NewFilterHistory(parent, &filter, &filter).Next(ctx)
	if err != nil {
		return false, err
	}

	if commit == nil {
		return false, nil
	}

	return true, nil
}
