package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := Root(cwd)
	if err != nil {
		return err
	}

	root := filepath.Join(path, DotDir)

	var cfg Config
	if err := ReadConfig(root, &cfg); err != nil {
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
			fmt.Printf("* ")
		}

		fmt.Println(b)
	}

	return WriteConfig(root, &cfg)
}
