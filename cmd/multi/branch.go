package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var branchCommand = &cli.Command{
	Action:    branchAction,
	Name:      "branch",
	Usage:     "List, create, or delete branches",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete branch",
		},
		&cli.BoolFlag{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create branch",
		},
	},
}

func branchAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := LoadConfig(cwd)
	if err != nil {
		return err
	}

	switch {
	case c.Bool("delete"):
		name := c.Args().Get(0)
		if name == "" {
			return errors.New("name cannot be empty")
		}

		if _, ok := config.Branches[name]; !ok {
			return errors.New("branch does not exists")
		}

		if name == config.Branch {
			return errors.New("cannot delete current branch")
		}

		delete(config.Branches, name)
	case c.Bool("create"):
		name := c.Args().Get(0)
		if name == "" {
			return errors.New("name cannot be empty")
		}

		if _, ok := config.Branches[name]; ok {
			return errors.New("branch already exists")
		}

		config.Branches[name] = config.Branches[config.Branch]
	}

	for branch := range config.Branches {
		if branch == config.Branch {
			fmt.Print("* ")
		}

		fmt.Println(branch)
	}

	return config.Save()
}
