package main

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"
)

var initCommand = &cli.Command{
	Action:    initAction,
	Name:      "init",
	Usage:     "Create a repo",
	ArgsUsage: "<name>",
}

func initAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := FindConfig(cwd); err == nil {
		return errors.New("repo already exists")
	}

	config := NewConfig(cwd, c.Args().Get(0))
	return config.Save()
}
