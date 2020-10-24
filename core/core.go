// Package core contains methods for interacting with Multiverse repositories.
package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-multiverse/ipfs"
)

var (
	// ErrInvalidFile is returned when an invalid file is encountered.
	ErrInvalidFile = errors.New("invalid file")
	// ErrInvalidRef is returned when a ref resolves to an invalid object.
	ErrInvalidRef = errors.New("invalid ref")
	// ErrMergeBase is returned when a merge base is not found.
	ErrMergeBase = errors.New("merge base not found")
	// ErrMergeAhead is returned when merge histories are equivalent.
	ErrMergeAhead = errors.New("local is ahead of remote")
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

// Worktree returns the current working tree files node.
func (c *Core) Worktree() (files.Node, error) {
	info, err := os.Stat(c.Config.Path)
	if err != nil {
		return nil, err
	}

	filter, err := files.NewFilter("", DefaultIgnore, true)
	if err != nil {
		return nil, err
	}

	return files.NewSerialFileWithFilter(c.Config.Path, filter, info)
}

// Checkout copies the tree of the commit with the given path to the local repo directory.
func (c *Core) Checkout(ctx context.Context, ref path.Path) error {
	p, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return err
	}

	if p.Cid().Type() != ipldmulti.CommitCodec {
		return ErrInvalidRef
	}

	node, err := c.Api.Unixfs().Get(ctx, path.Join(p, "tree"))
	if err != nil {
		return err
	}

	if err := writeNode(node, c.Config.Path); err != nil {
		return err
	}

	c.Config.Head = p.Root()
	return c.Config.Write()
}

// Commit records changes to the working directory.
func (c *Core) Commit(ctx context.Context, message string) (*ipldmulti.Commit, error) {
	tree, err := c.Worktree()
	if err != nil {
		return nil, err
	}

	p, err := c.Api.Unixfs().Add(ctx, tree)
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
		PeerID:   key.ID(),
		WorkTree: p.Root(),
	}

	if c.Config.Head.Defined() {
		commit.Parents = append(commit.Parents, c.Config.Head)
	}

	if err := c.Api.Dag().Pinning().Add(ctx, &commit); err != nil {
		return nil, err
	}

	c.Config.Head = commit.Cid()
	return &commit, c.Config.Write()
}

// Publish announces a new version to peers.
func (c *Core) Publish(ctx context.Context, name string, ref path.Path) (iface.IpnsEntry, error) {
	p, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return nil, err
	}

	if p.Cid().Type() != ipldmulti.CommitCodec {
		return nil, ErrInvalidRef
	}

	return c.Api.Name().Publish(ctx, ref, options.Name.Key(name))
}


// Diff prints the differences between the working directory and remote.
func (c *Core) Diff(ctx context.Context, ref path.Path) error {
	p, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return err
	}

	if p.Cid().Type() != ipldmulti.CommitCodec {
		return ErrInvalidRef
	}

	nodeA, err := c.Api.ResolveNode(ctx, path.Join(p, "tree"))
	if err != nil {
		return err
	}

	tree, err := c.Worktree()
	if err != nil {
		return err
	}

	p, err = c.Api.Unixfs().Add(ctx, tree)
	if err != nil {
		return err
	}

	nodeB, err := c.Api.ResolveNode(ctx, p)
	if err != nil {
		return err
	}

	diffs, err := dagutils.Diff(ctx, c.Api.Dag(), nodeA, nodeB)
	if err != nil {
		return err
	}

	for _, diff := range diffs {
		fmt.Println(diff)
	}

	return nil
}

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
