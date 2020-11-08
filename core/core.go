// Package core contains methods for interacting with Multiverse repositories.
package core

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-ipld-multiverse"
	"github.com/multiverse-vcs/go-multiverse/ipfs"
)

var (
	// ErrInvalidKey is returned when an invalid key is used.
	ErrInvalidKey = errors.New("invalid key")
	// ErrInvalidRef is returned when a ref resolves to an invalid object.
	ErrInvalidRef = errors.New("invalid ref")
	// ErrInvalidTree is returned when a commit has an invalid work tree.
	ErrInvalidTree = errors.New("invalid tree")
	// ErrInvalidFile is returned when an invalid file is encountered.
	ErrInvalidFile = errors.New("invalid file")
)

// Core contains core services.
type Core struct {
	api iface.CoreAPI
}

// CommitOptions are used set options when committing.
type CommitOptions struct {
	// Message describes the changes in the commit.
	Message string
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

// Commit creates a new commit containing a working tree and metadata.
func (c *Core) Commit(ctx context.Context, tree files.Node, opts *CommitOptions) (*ipldmulti.Commit, error) {
	key, err := c.api.Key().Self(ctx)
	if err != nil {
		return nil, err
	}

	p, err := c.api.Unixfs().Add(ctx, tree)
	if err != nil {
		return nil, err
	}

	commit := ipldmulti.Commit{
		Date:     time.Now(),
		Message:  opts.Message,
		Parents:  opts.Parents,
		PeerID:   key.ID(),
		WorkTree: p.Cid(),
	}

	adder := c.api.Dag().Pinning()
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

// ReadFile returns the contents of the file at the given path.
func (c *Core) ReadFile(ctx context.Context, path path.Path) (string, error) {
	node, err := c.api.Unixfs().Get(ctx, path)
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
func (c *Core) Status(ctx context.Context, ref path.Path, tree files.Node) ([]*dagutils.Change, error) {
	commit, err := c.Reference(ctx, ref)
	if err != nil {
		return nil, err
	}

	link, _, err := commit.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	p, err := c.api.Unixfs().Add(ctx, tree)
	if err != nil {
		return nil, err
	}

	nodeA, err := link.GetNode(ctx, c.api.Dag())
	if err != nil {
		return nil, err
	}

	nodeB, err := c.api.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(ctx, c.api.Dag(), nodeA, nodeB)
}

// Tree returns the work tree of the commit.
func (c *Core) Tree(ctx context.Context, commit *ipldmulti.Commit) (files.Node, error) {
	link, _, err := commit.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	tree, err := link.GetNode(ctx, c.api.Dag())
	if err != nil {
		return nil, err
	}

	node, ok := tree.(files.Node)
	if !ok {
		return nil, ErrInvalidTree
	}

	return node, nil
}
