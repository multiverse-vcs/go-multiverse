package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var initCommand = &cli.Command{
	Action:    initAction,
	Name:      "init",
	Usage:     "Initialize a new repo",
	ArgsUsage: "<cid>",
}

func initAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	repo := afero.NewBasePathFs(fs, cwd)
	root := filepath.Join(cwd, DotDir)

	if err := fs.Mkdir(root, 0755); err != nil {
		return err
	}

	node, err := node.NewNode(root)
	if err != nil {
		return err
	}

	cfg, err := DefaultConfig()
	if err != nil {
		return err
	}

	if c.NArg() == 0 {
		return WriteConfig(root, cfg)
	}

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	key, err := cfg.Key()
	if err != nil {
		return err
	}

	if err := node.Online(c.Context, key); err != nil {
		return err
	}

	fmt.Printf("fetching commit graph...\n")
	if err := merkledag.FetchGraph(c.Context, id, node.Dag); err != nil {
		return err
	}

	if err := core.Checkout(c.Context, repo, node.Dag, id); err != nil {
		return err
	}

	cfg.Index = id
	cfg.SetHead(id)

	return WriteConfig(root, cfg)
}
