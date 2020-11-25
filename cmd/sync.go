package cmd

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/p2p"
	"github.com/urfave/cli/v2"
)

// NewSyncCommand returns a new sync command.
func NewSyncCommand() *cli.Command {
	return &cli.Command{
		Name:      "sync",
		Usage:     "copy an existing repo",
		ArgsUsage: "[cid]",
		Action: func(c *cli.Context) error {
			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if c.NArg() < 1 {
				return cli.Exit("cid is required", 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if cfg.Head.Defined() {
				return cli.Exit("branch is not empty", 1)
			}

			id, err := cid.Parse(c.Args().Get(0))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.Online(c.Context); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("bootstrapping network...\n")
			p2p.Bootstrap(c.Context, store.Host)

			fmt.Printf("fetching commit graph...\n")
			if err := merkledag.FetchGraph(c.Context, id, store.Dag); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := core.Checkout(c.Context, store, id); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg.Head = id
			cfg.Base = id

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
