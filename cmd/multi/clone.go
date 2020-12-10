package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/urfave/cli/v2"
)

var cloneCommand = &cli.Command{
	Action:    cloneAction,
	Name:      "clone",
	Usage:     "Copy an existing repo",
	ArgsUsage: "<cid> <dir>",
}

func cloneAction(c *cli.Context) error {
	if c.NArg() < 2 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	root := filepath.Join(cwd, c.Args().Get(1), storage.DotDir)
	if err := fs.MkdirAll(root, 0755); err != nil {
		return err
	}

	store, err := openStore()
	if err != nil {
		return err
	}

	if err := store.Initialize(); err != nil {
		return err
	}

	cfg, err := store.ReadConfig()
	if err != nil {
		return err
	}

	if err := store.Online(c.Context); err != nil {
		return err
	}

	fmt.Printf("fetching commit graph...\n")
	if err := merkledag.FetchGraph(c.Context, id, store.Dag); err != nil {
		return err
	}

	if err := core.Checkout(c.Context, store, id); err != nil {
		return err
	}

	cfg.Index = id
	cfg.SetHead(id)

	return store.WriteConfig(cfg)
}
