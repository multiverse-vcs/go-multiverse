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
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multiaddr"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-multiverse/ipfs"
)

var (
	// DefaultIgnore contains default ignore rules.
	DefaultIgnore = []string{DefaultConfig, ".git"}
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
	p, err := c.Api.ResolvePath(ctx, path.Join(remote, "tree"))
	if err != nil {
		return err
	}

	tree, err := c.Api.Unixfs().Get(ctx, p)
	if err != nil {
		return err
	}

	if err := files.WriteTo(tree, c.Config.Path); err != nil {
		return err
	}

	c.Config.Head = p.Root()
	return c.Config.Write()
}

// Commit records changes to the working directory.
func (c *Core) Commit(ctx context.Context, message string) (*ipldmulti.Commit, error) {
	info, err := os.Stat(c.Config.Path)
	if err != nil {
		return nil, err
	}

	filter, err := files.NewFilter("", DefaultIgnore, true)
	if err != nil {
		return nil, err
	}

	tree, err := files.NewSerialFileWithFilter(c.Config.Path, filter, info)
	if err != nil {
		return nil, err
	}

	p, err := c.Api.Unixfs().Add(ctx, tree)
	if err != nil {
		return nil, err
	}

	commit := ipldmulti.Commit{Message: message, WorkTree: p.Root(), Date: time.Now()}
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

	fmt.Printf("Commit: %s\n", id.String())
	fmt.Printf("Date:   %s\n", commit.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
	fmt.Printf("\n%s\n\n", commit.Message)

	if len(commit.Parents) == 0 {
		return nil
	}

	return c.Log(ctx, commit.Parents[0])
}
