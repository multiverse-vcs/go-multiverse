package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// NewBranchCommand returns a new command.
func NewBranchCommand() *cli.Command {
	return &cli.Command{
		Name:  "branch",
		Usage: "add, remove, or list branches",
		Subcommands: []*cli.Command{
			NewBranchAddCommand(),
			NewBranchListCommand(),
			NewBranchRemoveCommand(),
		},
	}
}

// NewBranchAddCommand returns a new command.
func NewBranchAddCommand() *cli.Command {
	return &cli.Command{
		Name: "add",
		Usage: "create a branch",
		ArgsUsage: "<name>",
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

			if err := cfg.AddBranch(c.Args().Get(0), cfg.Head()); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}

// NewBranchListCommand returns a new command.
func NewBranchListCommand() *cli.Command {
	return &cli.Command{
		Name: "list",
		Aliases: []string{"ls"},
		Usage: "print branches",
		Action: func(c *cli.Context) error {
			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			for b, _ := range cfg.Branches {
				if b == cfg.Branch {
					fmt.Printf("%s*%s ", ColorGreen, ColorReset)
				}

				fmt.Println(b)
			}

			return nil
		},
	}
}

// NewBranchRemoveCommand returns a new command.
func NewBranchRemoveCommand() *cli.Command {
	return &cli.Command{
		Name: "remove",
		Aliases: []string{"rm"},
		Usage: "delete a branch",
		ArgsUsage: "<name>",
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

			if cfg.Branch == c.Args().Get(0) {
				return cli.Exit("cannot delete current branch", 1)
			}

			if err := cfg.DeleteBranch(c.Args().Get(0)); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
