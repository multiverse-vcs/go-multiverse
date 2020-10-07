// Package core contains methods for interacting with Multiverse repositories.
package core

import (
	"context"
	"fmt"
	"os"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multiaddr"
	"github.com/yondero/go-ipld-multiverse"
)

// DefaultIgnore contains default ignore rules.
var DefaultIgnore = []string{DefaultConfig, ".git"}

// Clone copies the tree of the commit with the given path.
func Clone(ctx context.Context, local string, remote string) (*Config, error) {
	addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		return nil, err
	}

	api, err := httpapi.NewApi(addr)
	if err != nil {
		return nil, err
	}

	p, err := api.ResolvePath(ctx, path.Join(path.New(remote), "tree"))
	if err != nil {
		return nil, err
	}

	f, err := api.Unixfs().Get(ctx, p)
	if err != nil {
		return nil, err
	}

	if err := files.WriteTo(f, local); err != nil {
		return nil, err
	}

	return InitConfig(local, p.Root())
}

// Commit records changes to the working directory.
func Commit(ctx context.Context, local string, message string) (*ipldmulti.Commit, error) {
	addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		return nil, err
	}

	api, err := httpapi.NewApi(addr)
	if err != nil {
		return nil, err
	}

	config, err := OpenConfig(local)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(config.Path)
	if err != nil {
		return nil, err
	}

	filter, err := files.NewFilter("", DefaultIgnore, true)
	if err != nil {
		return nil, err
	}

	tree, err := files.NewSerialFileWithFilter(config.Path, filter, info)
	if err != nil {
		return nil, err
	}

	p, err := api.Unixfs().Add(ctx, tree)
	if err != nil {
		return nil, err
	}

	c := ipldmulti.Commit{Message: message, WorkTree: p.Root()}
	if config.Head.Defined() {
		c.Parents = append(c.Parents, config.Head)
	}

	if err := api.Dag().Pinning().Add(ctx, &c); err != nil {
		return nil, err
	}

	config.Head = c.Cid()
	if err := config.Write(); err != nil {
		return nil, err
	}

	return &c, nil
}

// Log prints the commit history of the repo.
func Log(ctx context.Context, local string) error {
	addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		return err
	}

	api, err := httpapi.NewApi(addr)
	if err != nil {
		return err
	}

	config, err := OpenConfig(local)
	if err != nil {
		return err
	}

	id := config.Head
	for id.Defined() {
		node, err := api.Dag().Get(ctx, id)
		if err != nil {
			return err
		}

		c, ok := node.(*ipldmulti.Commit)
		if !ok {
			return nil
		}

		fmt.Printf("Commit: %s\n", c.Cid().String())
		fmt.Printf("Date:   %s\n", c.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
		fmt.Printf("\n%s\n\n", c.Message)

		if len(c.Parents) == 0 {
			return nil
		}

		id = c.Parents[0]
	}

	return nil
}
