package cmd

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

// NewCloneCommand returns a new command.
func NewCloneCommand() *cli.Command {
	return &cli.Command{
		Name:      "clone",
		Usage:     "Copy an existing repo",
		ArgsUsage: "<commit-cid> <dir>",
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

			root := filepath.Join(cwd, c.Args().Get(1))
			if err := os.Mkdir(root, 0755); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			store, err := storage.InitOsStore(root)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.Online(c.Context); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("fetching commit graph...\n")
			if err := merkledag.FetchGraph(c.Context, id, store.Dag); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := core.Checkout(c.Context, store, id); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg.Index = id
			cfg.SetHead(id)

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
