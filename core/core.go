// Package core contains methods for interacting with Multiverse repositories.
package core

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-ipld-multiverse"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/ipfs"
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

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{config.DefaultConfig}

// IgnoreFile is the name of the ignore file.
const IgnoreFile = ".multiverse.ignore"

// Core contains core services.
type Core struct {
	api iface.CoreAPI
}

// CommitOptions are used set options when committing.
type CommitOptions struct {
	// Message describes the changes in the commit.
	Message string
	// Pin specifies if the commit should be pinned.
	Pin bool
	// Parents are the ids of the parent commits.
	Parents []cid.Cid
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
func (c *Core) Checkout(ctx context.Context, commit *ipldmulti.Commit, root string) error {
	link, _, err := commit.ResolveLink([]string{"tree"})
	if err != nil {
		return err
	}

	tree, err := link.GetNode(ctx, c.api.Dag())
	if err != nil {
		return err
	}

	node, ok := tree.(files.Node)
	if !ok {
		return ErrInvalidTree
	}

	return WriteTree(node, root)
}

// Commit creates a new commit containing a working tree and metadata.
func (c *Core) Commit(ctx context.Context, tree cid.Cid, opts *CommitOptions) (*ipldmulti.Commit, error) {
	if !tree.Defined() {
		return nil, ErrInvalidTree
	}

	key, err := c.api.Key().Self(ctx)
	if err != nil {
		return nil, err
	}

	commit := ipldmulti.Commit{
		Date:     time.Now(),
		Message:  opts.Message,
		Parents:  opts.Parents,
		PeerID:   key.ID(),
		WorkTree: tree,
	}

	var adder format.NodeAdder = c.api.Dag()
	if opts.Pin {
		adder = c.api.Dag().Pinning()
	}

	if err := adder.Add(ctx, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// Diff returns the differences between two commit working trees.
func (c *Core) Diff(ctx context.Context, commitA, commitB *ipldmulti.Commit) ([]*dagutils.Change, error) {
	linkA, _, err := commitA.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	linkB, _, err := commitB.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	treeA, err := linkA.GetNode(ctx, c.api.Dag())
	if err != nil {
		return nil, err
	}

	treeB, err := linkB.GetNode(ctx, c.api.Dag())
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, c.api.Dag(), treeA, treeB)
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

// Reference resolves the commit from the given ref.
func (c *Core) Reference(ctx context.Context, ref path.Path) (*ipldmulti.Commit, error) {
	res, err := c.api.ResolvePath(ctx, ref)
	if err != nil {
		return nil, err
	}

	if res.Cid().Type() != ipldmulti.CommitCodec {
		return nil, ErrInvalidRef
	}

	node, err := c.api.Dag().Get(ctx, res.Cid())
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

	tree, err := c.WorkTree(ctx, root)
	if err != nil {
		return nil, err
	}

	nodeA, err := link.GetNode(ctx, c.api.Dag())
	if err != nil {
		return nil, err
	}

	nodeB, err := c.api.ResolveNode(ctx, tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, c.api.Dag(), nodeA, nodeB)
}

// WorkTree creates a tree from the given path and returns its cid.
func (c *Core) WorkTree(ctx context.Context, root string) (path.Resolved, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	ignore := filepath.Join(root, IgnoreFile)
	if _, err := os.Stat(ignore); err != nil {
		ignore = ""
	}

	filter, err := files.NewFilter(ignore, IgnoreRules, true)
	if err != nil {
		return nil, err
	}

	node, err := files.NewSerialFileWithFilter(root, filter, info)
	if err != nil {
		return nil, err
	}

	return c.api.Unixfs().Add(ctx, node)
}

// WriteTree writes the given node to the local repo root.
func WriteTree(node files.Node, root string) error {
	switch node := node.(type) {
	case *files.Symlink:
		return os.Symlink(node.Target, root)
	case files.File:
		b, err := ioutil.ReadAll(node)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(root, b, 0644)
	case files.Directory:
		if err := os.MkdirAll(root, 0777); err != nil {
			return err
		}

		entries := node.Entries()
		for entries.Next() {
			child := filepath.Join(root, entries.Name())
			if err := WriteTree(entries.Node(), child); err != nil {
				return err
			}
		}

		return entries.Err()
	default:
		return ErrInvalidFile
	}
}
