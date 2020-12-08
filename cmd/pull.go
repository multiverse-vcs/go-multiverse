package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/urfave/cli/v2"
)

// NewPullCommand returns a new command.
func NewPullCommand() *cli.Command {
	return &cli.Command{
		Name:      "pull",
		Usage:     "Merge changes from peers",
		ArgsUsage: "<commit-cid>",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			id, err := cid.Parse(c.Args().Get(0))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			store, err := storage.NewOsStore(cwd)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if cfg.Head() != cfg.Index {
				return cli.Exit("index is behind head", 1)
			}

			if err := store.Online(c.Context); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("fetching commit graph...\n")
			if err := merkledag.FetchGraph(c.Context, id, store.Dag); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			base, err := core.MergeBase(c.Context, store, cfg.Head(), id)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if base == id {
				return cli.Exit("local is ahead of remote", 1)
			}

			node, err := core.Merge(c.Context, store, cfg.Head(), base, id)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := core.Write(c.Context, store, "", node); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
