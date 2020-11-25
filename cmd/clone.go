package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/p2p"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/urfave/cli/v2"
)

// NewCloneCommand returns a new clone command.
func NewCloneCommand() *cli.Command {
	return &cli.Command{
		Name:      "clone",
		Usage:     "copy an existing repo",
		ArgsUsage: "<cid> <dir>",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return cli.Exit("missing required args", 1)
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

			dot := filepath.Join(root, storage.DotDir)
			if err := os.Mkdir(dot, 0755); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			store, err := storage.NewOsStore(root)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.Online(c.Context); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("bootstrapping network...\n")
			p2p.Bootstrap(c.Context, store.Host)

			if err := p2p.Discovery(c.Context, store.Host); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("fetching commit graph...\n")
			if err := merkledag.FetchGraph(c.Context, id, store.Dag); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := core.Checkout(c.Context, store, id); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg := config.Default()
			cfg.Head = id
			cfg.Base = id

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
