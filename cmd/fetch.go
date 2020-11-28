package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/p2p"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/urfave/cli/v2"
)

// NewFetchCommand returns a new command.
func NewFetchCommand() *cli.Command {
	return &cli.Command{
		Name:      "fetch",
		Usage:     "copy changes from peers",
		ArgsUsage: "<commit-cid> <branch-name>",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
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

			name := c.Args().Get(1)
			if b, ok := cfg.Branches[name]; ok && b.Head.Defined() {
				return cli.Exit("branch exists and is not empty", 1)
			}

			if err := store.Online(c.Context); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("bootstrapping network...\n")
			p2p.Bootstrap(c.Context, store.Host)

			if err := p2p.Discovery(c.Context, store.Host); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.Router.Bootstrap(c.Context); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("fetching commit graph...\n")
			if err := merkledag.FetchGraph(c.Context, id, store.Dag); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg.Branches[name] = &config.Branch{Head: id}
			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
