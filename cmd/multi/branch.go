package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var branchCommand = &cli.Command{
	Action:    branchAction,
	Name:      "branch",
	Usage:     "Add, remove, or list branches",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Value:   false,
			Usage:   "Delete branch",
		},
	},
}

func branchAction(c *cli.Context) error {
	store, err := openStore()
	if err != nil {
		return err
	}

	cfg, err := store.ReadConfig()
	if err != nil {
		return err
	}

	switch {
	case c.NArg() > 0 && c.Bool("delete"):
		if cfg.Branch == c.Args().Get(0) {
			return errors.New("cannot delete current branch")
		}

		if err := cfg.DeleteBranch(c.Args().Get(0)); err != nil {
			return err
		}
	case c.NArg() > 0:
		if err := cfg.AddBranch(c.Args().Get(0), cfg.Head()); err != nil {
			return err
		}
	}

	for b := range cfg.Branches {
		if b == cfg.Branch {
			fmt.Printf("%s*%s ", ColorGreen, ColorReset)
		}

		fmt.Println(b)
	}

	return store.WriteConfig(cfg)
}
