package cmd

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/urfave/cli/v2"
)

// NewSwitchCommand returns a new command.
func NewSwitchCommand() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Usage:     "change branches",
		ArgsUsage: "<branch-name>",
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

			name := c.Args().Get(0)
			if cfg.Branch == name {
				return cli.Exit("already on branch", 1)
			}

			branch, err := cfg.GetBranch(name)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Println("stashing changes...")
			stash, err := core.Worktree(c.Context, store)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			var id cid.Cid = branch.Head
			if branch.Stash.Defined() {
				id = branch.Stash
			}

			fmt.Println("checking out branch...")
			if err := core.Checkout(c.Context, store, id); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg.SetStash(stash.Cid())
			cfg.Branch = name
			cfg.Index = id

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
