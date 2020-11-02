package core

import (
	"context"
	"io/ioutil"
	"sort"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-xdiff"
)

// Merge combines the repo histories of two commits into a single tree and returns the root node.
func (c *Core) Merge(ctx context.Context, local, remote *ipldmulti.Commit) (*merkledag.ProtoNode, error) {
	base, err := c.MergeBase(ctx, local.Cid(), remote.Cid())
	if err != nil {
		return nil, err
	}

	if base.Cid() == local.Cid() {
		return nil, ErrMergeBehind
	}

	changes, err := c.MergeDiffs(ctx, local, remote, base)
	if err != nil {
		return nil, err
	}

	link, _, err := base.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	tree, err := link.GetNode(ctx, c.Api.Dag())
	if err != nil {
		return nil, err
	}

	proto, ok := tree.(*merkledag.ProtoNode)
	if !ok {
		return nil, ErrInvalidRef
	}

	return dagutils.ApplyChange(ctx, c.Api.Dag(), proto, changes)
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

// MergeConflict creates a new change that resolves the conflicting changes.
func (c *Core) MergeConflict(ctx context.Context, conflict dagutils.Conflict) (*dagutils.Change, error) {
	if conflict.A.Type == dagutils.Remove && conflict.B.Type == dagutils.Remove {
		return conflict.A, nil
	}

	merge, err := c.MergeFiles(ctx, conflict.A.Before, conflict.A.After, conflict.B.After)
	if err != nil {
		return nil, err
	}

	change := dagutils.Change{
		Type:   dagutils.Mod,
		Path:   conflict.A.Path,
		Before: conflict.A.Before,
		After:  merge.Cid(),
	}

	if conflict.A.Type == dagutils.Add && conflict.B.Type == dagutils.Add {
		change.Type = dagutils.Add
	}

	return &change, nil
}

// MergeDiffs merges the changes from local and remote using base as a common ancestor.
func (c *Core) MergeDiffs(ctx context.Context, local, remote, base *ipldmulti.Commit) ([]*dagutils.Change, error) {
	ours, err := c.DiffWorkTrees(ctx, base, local)
	if err != nil {
		return nil, err
	}

	theirs, err := c.DiffWorkTrees(ctx, base, remote)
	if err != nil {
		return nil, err
	}

	changes, conflicts := dagutils.MergeDiffs(ours, theirs)
	for _, conflict := range conflicts {
		change, err := c.MergeConflict(ctx, conflict)
		if err != nil {
			return nil, err
		}

		changes = append(changes, change)
	}

	return changes, nil
}

// MergeFiles creates a new file by performing a three way merge using base, local, and remote.
func (c *Core) MergeFiles(ctx context.Context, base, local, remote cid.Cid) (path.Resolved, error) {
	original, err := c.readChange(ctx, base)
	if err != nil {
		return nil, err
	}

	ours, err := c.readChange(ctx, local)
	if err != nil {
		return nil, err
	}

	theirs, err := c.readChange(ctx, remote)
	if err != nil {
		return nil, err
	}

	merge, err := xdiff.Merge(original, ours, theirs, &xdiff.DefaultMergeOptions)
	if err != nil {
		return nil, err
	}

	return c.Api.Unixfs().Add(ctx, files.NewBytesFile([]byte(merge)))
}

// readChange returns a the contents of a change.
func (c *Core) readChange(ctx context.Context, id cid.Cid) (string, error) {
	if !id.Defined() {
		return "", nil
	}

	node, err := c.Api.Unixfs().Get(ctx, path.IpfsPath(id))
	if err != nil {
		return "", err
	}

	file, ok := node.(files.File)
	if !ok {
		return "", ErrInvalidFile
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
