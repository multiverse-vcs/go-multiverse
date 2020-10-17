// Package core contains methods for interacting with Multiverse repositories.
package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
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
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo cannot be found.
	ErrRepoNotFound = errors.New("repo not found")
	// ErrInvalidRef is returned when a ref resolves to an invalid object.
	ErrInvalidRef = errors.New("ref is not a multiverse object")
	// ErrInvalidFile is returned when an invalid file is encountered.
	ErrInvalidFile = errors.New("invalid file type")
)

// Core contains config and core services.
type Core struct {
	// Config is the local repo config.
	Config *Config
	// Api is an IPFS core api.
	Api iface.CoreAPI
}

// DefaultIgnore contains default ignore rules.
var DefaultIgnore = []string{DefaultConfig, ".git"}

// NewCore returns a new core api.
func NewCore(ctx context.Context, config *Config) (*Core, error) {
	api, err := ipfs.NewApi(ctx)
	if err != nil {
		return nil, err
	}

	return &Core{config, api}, nil
}

// Tree returns the current working tree files node.
func (c *Core) Tree() (files.Node, error) {
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

// Commit records changes to the working directory.
func (c *Core) Commit(ctx context.Context, message string) (*ipldmulti.Commit, error) {
	tree, err := c.Tree()
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

// Log prints the commit history of the repo.
func (c *Core) Log(ctx context.Context, id cid.Cid) error {
	if !id.Defined() {
		return nil
	}

	node, err := c.Api.Dag().Get(ctx, id)
	if err != nil {
		return err
	}

	commit, ok := node.(*ipldmulti.Commit)
	if !ok {
		return nil
	}

	color.Yellow.Printf("commit %s\n", id.String())
	fmt.Printf("Peer: %s\n", commit.PeerID.String())
	fmt.Printf("Date: %s\n", commit.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
	fmt.Printf("\n\t%s\n\n", commit.Message)

	if len(commit.Parents) == 0 {
		return nil
	}

	return c.Log(ctx, commit.Parents[0])
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

	tree, err := c.Tree()
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
