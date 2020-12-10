package main

import (
	"errors"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/urfave/cli/v2"
)

var pullCommand = &cli.Command{
	Action:    pullAction,
	Name:      "pull",
	Usage:     "Merge changes into the current branch",
	ArgsUsage: "<cid>",
}

func pullAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	store, err := openStore()
	if err != nil {
		return err
	}

	cfg, err := store.ReadConfig()
	if err != nil {
		return err
	}

	if cfg.Head() != cfg.Index {
		return errors.New("index is behind head")
	}

	if err := store.Online(c.Context); err != nil {
		return err
	}

	fmt.Printf("fetching commit graph...\n")
	if err := merkledag.FetchGraph(c.Context, id, store.Dag); err != nil {
		return err
	}

	base, err := core.MergeBase(c.Context, store, cfg.Head(), id)
	if err != nil {
		return err
	}

	if base == id {
		return errors.New("local is ahead of remote")
	}

	node, err := core.Merge(c.Context, store, cfg.Head(), base, id)
	if err != nil {
		return err
	}

	if err := core.Write(c.Context, store, "", node); err != nil {
		return err
	}

	return nil
}
