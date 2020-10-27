// Package core contains methods for interacting with Multiverse repositories.
package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/go-unixfs/file"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-multiverse/ipfs"
)

var (
	// ErrInvalidFile is returned when an invalid file is encountered.
	ErrInvalidFile = errors.New("invalid file")
	// ErrInvalidKey is returned when an invalid key is used.
	ErrInvalidKey = errors.New("invalid key")
	// ErrInvalidRef is returned when a ref resolves to an invalid object.
	ErrInvalidRef = errors.New("invalid ref")
	// ErrMergeBase is returned when a merge base is not found.
	ErrMergeBase = errors.New("merge base not found")
	// ErrMergeAhead is returned when merge histories are equivalent.
	ErrMergeAhead = errors.New("local is ahead of remote")
	// ErrNoChanges is returned when there are no changes to commit.
	ErrNoChanges = errors.New("no changes to commit")
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo cannot be found.
	ErrRepoNotFound = errors.New("repo not found")
)

// DefaultIgnore contains default ignore rules.
var DefaultIgnore = []string{DefaultConfig, ".git"}

// Core contains config and core services.
type Core struct {
	// Config is the local repo config.
	Config *Config
	// Api is an IPFS core api.
	Api iface.CoreAPI
}

// NewCore returns a new core api.
func NewCore(ctx context.Context, config *Config) (*Core, error) {
	api, err := ipfs.NewApi(ctx)
	if err != nil {
		return nil, err
	}

	return &Core{config, api}, nil
}

// Checkout copies the tree of the commit with the given path to the local repo directory.
func (c *Core) Checkout(ctx context.Context, ref path.Path) (*ipldmulti.Commit, error) {
	node, err := c.Api.ResolveNode(ctx, ref)
	if err != nil {
		return nil, err
	}

	commit, ok := node.(*ipldmulti.Commit)
	if !ok {
		return nil, ErrInvalidRef
	}

	tree, err := c.Api.Unixfs().Get(ctx, path.Join(ref, "tree"))
	if err != nil {
		return nil, err
	}

	if err := writeNode(tree, c.Config.Path); err != nil {
		return nil, err
	}

	return commit, nil
}

// Commit records changes to the working directory.
func (c *Core) Commit(ctx context.Context, message string, parents ...cid.Cid) (*ipldmulti.Commit, error) {
	tree, err := c.WorkTree(ctx)
	if err != nil {
		return nil, err
	}

	key, err := c.Api.Key().Self(ctx)
	if err != nil {
		return nil, err
	}

	commit := ipldmulti.Commit{
		Date:     time.Now(),
		Message:  message,
		Parents:  parents,
		PeerID:   key.ID(),
		WorkTree: tree.Cid(),
	}

	if err := c.Api.Dag().Pinning().Add(ctx, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// Diff returns the differences between two commit working trees.
func (c *Core) Diff(ctx context.Context, refA, refB path.Path) ([]*dagutils.Change, error) {
	nodeA, err := c.Api.ResolveNode(ctx, path.Join(refA, "tree"))
	if err != nil {
		return nil, err
	}

	nodeB, err := c.Api.ResolveNode(ctx, path.Join(refB, "tree"))
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, c.Api.Dag(), nodeA, nodeB)
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

// Merge combines the repo histories of the local and remote commits.
func (c *Core) Merge(ctx context.Context, ref path.Path) (*ipldmulti.Commit, error) {
	p, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return nil, err
	}

	if p.Cid().Type() != ipldmulti.CommitCodec {
		return nil, ErrInvalidRef
	}

	base, err := c.MergeBase(ctx, c.Config.Head, p.Cid())
	if err != nil {
		return nil, err
	}

	if base.Cid().Equals(c.Config.Head) {
		return c.Checkout(ctx, ref)
	}

	local, err := c.Diff(ctx, base, path.IpfsPath(c.Config.Head))
	if err != nil {
		return nil, err
	}

	remote, err := c.Diff(ctx, base, p)
	if err != nil {
		return nil, err
	}

	changes, conflicts := dagutils.MergeDiffs(local, remote)
	if err := c.MergeChanges(ctx, base, changes); err != nil {
		return nil, err
	}

	if err := c.MergeConflicts(ctx, conflicts); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("merge %s", p.Cid().String())
	return c.Commit(ctx, message, c.Config.Head, p.Cid())
}

// MergeBase returns the best merge base for local and remote.
func (c *Core) MergeBase(ctx context.Context, local, remote cid.Cid) (path.Resolved, error) {
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

	iter := c.NewFilterHistory(remote, &filter, &filter)
	if err := iter.ForEach(ctx, callback); err != nil {
		return nil, err
	}

	if len(bases) == 0 {
		return nil, ErrMergeBase
	}

	sort.Slice(bases, func(i, j int) bool {
		return bases[i].Date.After(bases[j].Date)
	})

	return path.IpfsPath(bases[0].Cid()), nil
}

// MergeChanges merges the changes into the working tree and writes them to the local repo.
func (c *Core) MergeChanges(ctx context.Context, ref path.Path, changes []*dagutils.Change) error {
	base, err := c.Api.ResolveNode(ctx, path.Join(ref, "tree"))
	if err != nil {
		return err
	}

	proto, ok := base.(*merkledag.ProtoNode)
	if !ok {
		return ErrInvalidRef
	}

	node, err := dagutils.ApplyChange(ctx, c.Api.Dag(), proto, changes)
	if err != nil {
		return err
	}

	tree, err := unixfile.NewUnixfsFile(ctx, c.Api.Dag(), node)
	if err != nil {
		return err
	}

	return writeNode(tree, c.Config.Path)
}

// MergeConflicts writes the conflicts to the local repo so they can be resolved manually.
func (c *Core) MergeConflicts(ctx context.Context, conflicts []dagutils.Conflict) error {
	return nil
}

// Publish announces a new version to peers.
func (c *Core) Publish(ctx context.Context, key string, ref path.Path) (iface.IpnsEntry, error) {
	if key == "self" {
		return nil, ErrInvalidKey
	}

	p, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return nil, err
	}

	if p.Cid().Type() != ipldmulti.CommitCodec {
		return nil, ErrInvalidRef
	}

	return c.Api.Name().Publish(ctx, ref, options.Name.Key(key))
}

// Status returns changes between local repo and head.
func (c *Core) Status(ctx context.Context) ([]*dagutils.Change, error) {
	head := path.IpfsPath(c.Config.Head)

	nodeA, err := c.Api.ResolveNode(ctx, path.Join(head, "tree"))
	if err != nil {
		return nil, err
	}

	tree, err := c.WorkTree(ctx)
	if err != nil {
		return nil, err
	}

	nodeB, err := c.Api.ResolveNode(ctx, tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, c.Api.Dag(), nodeA, nodeB)
}

// Worktree adds the local repo changes and returns its cid.
func (c *Core) WorkTree(ctx context.Context) (path.Resolved, error) {
	info, err := os.Stat(c.Config.Path)
	if err != nil {
		return nil, err
	}

	filter, err := files.NewFilter("", DefaultIgnore, true)
	if err != nil {
		return nil, err
	}

	node, err := files.NewSerialFileWithFilter(c.Config.Path, filter, info)
	if err != nil {
		return nil, err
	}

	return c.Api.Unixfs().Add(ctx, node)
}
