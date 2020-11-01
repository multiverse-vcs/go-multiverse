// Package core contains methods for interacting with Multiverse repositories.
package core

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-multiverse/config"
	"github.com/yondero/go-multiverse/ipfs"
)

var (
	// ErrInvalidFile is returned when an invalid file is encountered.
	ErrInvalidFile = errors.New("invalid file")
	// ErrInvalidKey is returned when an invalid key is used.
	ErrInvalidKey = errors.New("invalid key")
	// ErrInvalidRef is returned when a ref resolves to an invalid object.
	ErrInvalidRef = errors.New("invalid ref")
	// ErrInvalidTree is returned when a commit has an invalid work tree.
	ErrInvalidTree = errors.New("work tree is invalid")
	// ErrMergeBase is returned when a merge base is not found.
	ErrMergeBase = errors.New("merge base not found")
	// ErrMergeAhead is returned when local contains remote changes.
	ErrMergeAhead = errors.New("local is ahead of remote")
	// ErrMergeBehind is returned when  remote contains local changes.
	ErrMergeBehind = errors.New("local is behind remote")
	// ErrNoChanges is returned when there are no changes to commit.
	ErrNoChanges = errors.New("no changes to commit")
)

// DefaultIgnore contains default ignore rules.
var DefaultIgnore = []string{config.DefaultConfig, ".git"}

// Core contains core services.
type Core struct {
	// Api is an IPFS core api.
	Api iface.CoreAPI
}

// CommitOptions are used set options when committing.
type CommitOptions struct {
	// Message describes the changes in the commit.
	Message  string
	// Pin specifies if the commit should be pinned.
	Pin      bool
	// Parents are the ids of the parent commits.
	Parents  []cid.Cid
	// WorkTree is the id of the changes in the commit.
	WorkTree cid.Cid
}

// NewCore returns a new core api.
func NewCore(ctx context.Context) (*Core, error) {
	api, err := ipfs.NewApi(ctx)
	if err != nil {
		return nil, err
	}

	return &Core{api}, nil
}

// Checkout copies the tree of the commit with the given ref to the local repo directory.
func (c *Core) Checkout(ctx context.Context, ref path.Path, root string) (*ipldmulti.Commit, error) {
	commit, err := c.Reference(ctx, ref)
	if err != nil {
		return nil, err
	}

	link, _, err := commit.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	tree, err := link.GetNode(ctx, c.Api.Dag())
	if err != nil {
		return nil, err
	}

	node, ok := tree.(files.Node)
	if !ok {
		return nil, ErrInvalidTree
	}

	if err := writeNode(node, root); err != nil {
		return nil, err
	}

	return commit, nil
}

// Commit creates a new commit containing a working tree and metadata.
func (c *Core) Commit(ctx context.Context, opts *CommitOptions) (*ipldmulti.Commit, error) {
	if !opts.WorkTree.Defined() {
		return nil, ErrInvalidTree
	}

	key, err := c.Api.Key().Self(ctx)
	if err != nil {
		return nil, err
	}

	commit := ipldmulti.Commit{
		Date:     time.Now(),
		Message:  opts.Message,
		Parents:  opts.Parents,
		PeerID:   key.ID(),
		WorkTree: opts.WorkTree,
	}

	var adder format.NodeAdder = c.Api.Dag()
	if opts.Pin {
		adder = c.Api.Dag().Pinning()
	}

	if err := adder.Add(ctx, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// DiffWorkTrees returns the differences between two commit working trees.
func (c *Core) DiffWorkTrees(ctx context.Context, commitA, commitB *ipldmulti.Commit) ([]*dagutils.Change, error) {
	linkA, _, err := commitA.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	linkB, _, err := commitB.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	treeA, err := linkA.GetNode(ctx, c.Api.Dag())
	if err != nil {
		return nil, err
	}

	treeB, err := linkB.GetNode(ctx, c.Api.Dag())
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, c.Api.Dag(), treeA, treeB)
}

// IsAncestor checks if child is an ancestor of parent.
func (c *Core) IsAncestor(ctx context.Context, child, parent cid.Cid) (bool, error) {
	var filter HistoryFilter = func(commit *ipldmulti.Commit) bool {
		return commit.Cid().Equals(child)
	}

	commit, err := c.NewHistory(parent).WithFilter(&filter, &filter).Next(ctx)
	if err != nil {
		return false, err
	}

	if commit == nil {
		return false, nil
	}

	return true, nil
}

// Publish announces a new version to peers.
func (c *Core) Publish(ctx context.Context, key string, ref path.Path) (iface.IpnsEntry, error) {
	if key == "self" {
		return nil, ErrInvalidKey
	}

	_, err := c.Reference(ctx, ref)
	if err != nil {
		return nil, err
	}

	return c.Api.Name().Publish(ctx, ref, options.Name.Key(key))
}

// Reference resolves the commit from the given ref.
func (c *Core) Reference(ctx context.Context, ref path.Path) (*ipldmulti.Commit, error) {
	res, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return nil, err
	}

	if res.Cid().Type() != ipldmulti.CommitCodec {
		return nil, ErrInvalidRef
	}

	node, err := c.Api.Dag().Get(ctx, res.Cid())
	if err != nil {
		return nil, err
	}

	commit, ok := node.(*ipldmulti.Commit)
	if !ok {
		return nil, ErrInvalidRef
	}

	return commit, nil
}

// Status returns changes between local repo and head.
func (c *Core) Status(ctx context.Context, ref path.Path, root string) ([]*dagutils.Change, error) {
	commit, err := c.Reference(ctx, ref)
	if err != nil {
		return nil, err
	}

	link, _, err := commit.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	treeA, err := c.WorkTree(ctx, root)
	if err != nil {
		return nil, err
	}

	nodeA, err := c.Api.ResolveNode(ctx, treeA)
	if err != nil {
		return nil, err
	}

	nodeB, err := link.GetNode(ctx, c.Api.Dag())
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, c.Api.Dag(), nodeA, nodeB)
}

// WorkTree creates a tree from the given path and returns its cid.
func (c *Core) WorkTree(ctx context.Context, root string) (path.Resolved, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	filter, err := files.NewFilter("", DefaultIgnore, true)
	if err != nil {
		return nil, err
	}

	node, err := files.NewSerialFileWithFilter(root, filter, info)
	if err != nil {
		return nil, err
	}

	return c.Api.Unixfs().Add(ctx, node)
}
