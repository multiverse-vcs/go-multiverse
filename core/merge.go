package core

import (
	"context"
	"errors"
	"sort"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/go-unixfs/file"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-ipld-multiverse"
	"github.com/multiverse-vcs/go-xdiff"
)

var (
	// ErrMergeBase is returned when a merge base is not found.
	ErrMergeBase = errors.New("merge base not found")
	// ErrMergeAhead is returned when local contains remote changes.
	ErrMergeAhead = errors.New("local is ahead of remote")
	// ErrMergeBehind is returned when  remote contains local changes.
	ErrMergeBehind = errors.New("local is behind remote")
)

// Merge combines the repo histories of two commits into a single tree and returns the root node.
func (c *Core) Merge(ctx context.Context, local, remote *ipldmulti.Commit) (files.Node, error) {
	base, err := c.MergeBase(ctx, local.Cid(), remote.Cid())
	if err != nil {
		return nil, err
	}

	if base.Cid() == local.Cid() {
		return nil, ErrMergeBehind
	}

	changes, err := c.MergeChanges(ctx, base, local, remote)
	if err != nil {
		return nil, err
	}

	link, _, err := base.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	tree, err := link.GetNode(ctx, c.api.Dag())
	if err != nil {
		return nil, err
	}

	proto, ok := tree.(*merkledag.ProtoNode)
	if !ok {
		return nil, ErrInvalidRef
	}

	merge, err := dagutils.ApplyChange(ctx, c.api.Dag(), proto, changes)
	if err != nil {
		return nil, err
	}

	return unixfile.NewUnixfsFile(ctx, c.api.Dag(), merge)
}

// MergeBase returns the best merge base for local and remote.
func (c *Core) MergeBase(ctx context.Context, local, remote cid.Cid) (*ipldmulti.Commit, error) {
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

	var bases []*ipldmulti.Commit
	var callback HistoryCallback = func(commit *ipldmulti.Commit) error {
		bases = append(bases, commit)
		return nil
	}

	iter := c.NewHistory(remote).WithFilter(&filter, &filter)
	if err := iter.ForEach(ctx, callback); err != nil {
		return nil, err
	}

	if len(bases) == 0 {
		return nil, ErrMergeBase
	}

	// TODO find least common ancestor
	sort.Slice(bases, func(i, j int) bool {
		return bases[i].Date.After(bases[j].Date)
	})

	return bases[0], nil
}

// MergeChanges merges the changes from local and remote using base as a common ancestor.
func (c *Core) MergeChanges(ctx context.Context, base, local, remote *ipldmulti.Commit) ([]*dagutils.Change, error) {
	ours, err := c.Diff(ctx, base, local)
	if err != nil {
		return nil, err
	}

	theirs, err := c.Diff(ctx, base, remote)
	if err != nil {
		return nil, err
	}

	changes, conflicts := dagutils.MergeDiffs(ours, theirs)

	// resolve conflicts and append to changes
	for _, conflict := range conflicts {
		change, err := c.MergeConflict(ctx, conflict)
		if err != nil {
			return nil, err
		}

		changes = append(changes, change)
	}

	return changes, nil
}

// MergeConflict creates a new change that resolves the conflicting changes.
func (c *Core) MergeConflict(ctx context.Context, conflict dagutils.Conflict) (*dagutils.Change, error) {
	if conflict.A.Type == dagutils.Remove && conflict.B.Type == dagutils.Remove {
		return conflict.A, nil
	}

	merge, err := c.MergeFiles(ctx, conflict.A.Before, conflict.A.After, conflict.B.After)
	if err != nil {
		return nil, err
	}

	p, err := c.api.Unixfs().Add(ctx, files.NewBytesFile([]byte(merge)))
	if err != nil {
		return nil, err
	}

	change := dagutils.Change{
		Type:   dagutils.Mod,
		Path:   conflict.A.Path,
		Before: conflict.A.Before,
		After:  p.Cid(),
	}

	if conflict.A.Type == dagutils.Add && conflict.B.Type == dagutils.Add {
		change.Type = dagutils.Add
	}

	return &change, nil
}

// MergeFiles merges the contents of local and remote into base.
func (c *Core) MergeFiles(ctx context.Context, base, local, remote cid.Cid) (string, error) {
	var err error
	var original, ours, theirs string

	if base.Defined() {
		original, err = c.ReadFile(ctx, path.IpfsPath(base))
	}

	if err != nil {
		return "", err
	}

	if local.Defined() {
		ours, err = c.ReadFile(ctx, path.IpfsPath(local))
	}

	if err != nil {
		return "", err
	}

	if remote.Defined() {
		theirs, err = c.ReadFile(ctx, path.IpfsPath(remote))
	}

	if err != nil {
		return "", err
	}

	return xdiff.Merge(original, ours, theirs, &xdiff.DefaultMergeOptions)
}
