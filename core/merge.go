package core

import (
	"bytes"
	"context"
	"io"
	"sort"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-xdiff"
)

// Merge combines the repo histories of two commits.
func (c *Core) Merge(ctx context.Context, ref path.Path, message string) (*ipldmulti.Commit, error) {
	local := path.IpfsPath(c.Config.Head)

	remote, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return nil, err
	}

	if remote.Cid().Type() != ipldmulti.CommitCodec {
		return nil, ErrInvalidRef
	}

	base, err := c.MergeBase(ctx, local.Cid(), remote.Cid())
	if err != nil {
		return nil, err
	}

	if base.Cid() == local.Cid() {
		return c.Checkout(ctx, ref)
	}

	ours, err := c.Diff(ctx, base, local)
	if err != nil {
		return nil, err
	}

	theirs, err := c.Diff(ctx, base, remote)
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

	tree, err := c.Api.ResolveNode(ctx, path.Join(ref, "tree"))
	if err != nil {
		return nil, err
	}

	proto, ok := tree.(*merkledag.ProtoNode)
	if !ok {
		return nil, ErrInvalidRef
	}

	merge, err := dagutils.ApplyChange(ctx, c.Api.Dag(), proto, changes)
	if err != nil {
		return nil, err
	}

	return c.Commit(ctx, path.IpfsPath(merge.Cid()), message, local.Cid(), remote.Cid())
}

// MergeBase returns the best merge base for local and remote.
func (c *Core) MergeBase(ctx context.Context, local, remote cid.Cid) (path.Resolved, error) {
	history, err := c.NewHistory(local).Flatten(ctx)
	if err != nil {
		return nil, err
	}

	// local is ahead of remote
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

	iter := c.NewFilterHistory(remote, &filter, &filter)
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

	return path.IpfsPath(bases[0].Cid()), nil
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

// MergeFiles creates a new file by performing a three way merge using base, local, and remote.
func (c *Core) MergeFiles(ctx context.Context, base, local, remote cid.Cid) (path.Resolved, error) {
	original, err := c.nodeReader(ctx, base)
	if err != nil {
		return nil, err
	}

	ours, err := c.nodeReader(ctx, local)
	if err != nil {
		return nil, err
	}

	theirs, err := c.nodeReader(ctx, remote)
	if err != nil {
		return nil, err
	}

	merge, err := xdiff.Merge(original, ours, theirs, &xdiff.DefaultMergeOptions)
	if err != nil {
		return nil, err
	}

	return c.Api.Unixfs().Add(ctx, files.NewBytesFile([]byte(merge)))
}

// nodeReader returns a reader for the contents of the file node with the given id.
func (c *Core) nodeReader(ctx context.Context, id cid.Cid) (io.Reader, error) {
	if !id.Defined() {
		return bytes.NewReader(nil), nil
	}

	node, err := c.Api.Unixfs().Get(ctx, path.IpfsPath(id))
	if err != nil {
		return nil, err
	}

	file, ok := node.(files.File)
	if !ok {
		return nil, ErrInvalidFile
	}

	return file, nil
}
