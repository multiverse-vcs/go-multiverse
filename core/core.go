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
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multiaddr"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-multiverse/file"
	"github.com/yondero/go-multiverse/ipfs"
)

var (
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo cannot be found.
	ErrRepoNotFound = errors.New("repo not found")
	// ErrIpfsApi is returned when a connection to the local ipfs node fails.
	ErrIpfsApi = errors.New("failed to connect to local ipfs node")
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
func NewCore(config *Config) (*Core, error) {
	addr, err := multiaddr.NewMultiaddr(ipfs.CommandsApiAddress)
	if err != nil {
		return nil, err
	}

	api, err := httpapi.NewApi(addr)
	if err != nil {
		return nil, ErrIpfsApi
	}

	return &Core{config, api}, nil
}

// Checkout copies the tree of the commit with the given path to the local repo directory.
func (c *Core) Checkout(ctx context.Context, remote path.Path) error {
	p, err := c.Api.ResolvePath(ctx, remote)
	if err != nil {
		return err
	}

	node, err := c.Api.Unixfs().Get(ctx, path.Join(p, "tree"))
	if err != nil {
		return err
	}

	entries := node.(files.Directory).Entries()
	if err := file.WriteEntries(entries, c.Config.Path); err != nil {
		return err
	}

	c.Config.Head = p.Root()
	return c.Config.Write()
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
		Date: time.Now(),
		Message: message,
		PeerID: key.ID(),
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

// Diff prints differences between two commits.
func (c *Core) Diff(ctx context.Context) error {
	head := path.IpfsPath(c.Config.Head)

	nodeA, err := c.Api.Unixfs().Get(ctx, path.Join(head, "tree"))
	if err != nil {
		return err
	}

	nodeB, err := c.Tree()
	if err != nil {
		return err
	}

	diffs, err := file.Diff(nodeA, nodeB)
	if err != nil {
		return err
	}

	for _, diff := range diffs {
		patch, err := diff.Patch()
		if err != nil {
			return err
		}

		color.Bold.Println(diff.Path)
		fmt.Println(patch)
	}

	return nil
}

// Status prints the differences between the working directory and remote.
func (c *Core) Status(ctx context.Context) error {
	head := path.IpfsPath(c.Config.Head)

	nodeA, err := c.Api.Unixfs().Get(ctx, path.Join(head, "tree"))
	if err != nil {
		return err
	}

	nodeB, err := c.Tree()
	if err != nil {
		return err
	}

	diffs, err := file.Diff(nodeA, nodeB)
	if err != nil {
		return err
	}

	for _, diff := range diffs {
		fmt.Println(diff)
	}

	return nil
}
