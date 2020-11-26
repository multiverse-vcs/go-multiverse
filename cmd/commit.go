package cmd

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/urfave/cli/v2"
)

// NewCommitCommand returns a new commit command.
func NewCommitCommand() *cli.Command {
	return &cli.Command{
		Name:      "commit",
		Usage:     "record repo changes",
		ArgsUsage: "<message>",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return cli.Exit("missing required args", 1)
			}

			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			var parents []cid.Cid
			if cfg.Head.Defined() {
				parents = append(parents, cfg.Head)
			}

			id, err := core.Commit(c.Context, store, c.Args().Get(0), parents...)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg.Head = id
			cfg.Base = id

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("%s%s%s\n", ColorYellow, id.String(), ColorReset)
			return nil
		},
	}
}
