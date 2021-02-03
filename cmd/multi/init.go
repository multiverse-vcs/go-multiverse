package main

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var initCommand = &cli.Command{
	Action:    initAction,
	Name:      "init",
	Usage:     "Create a repo",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Default branch",
			Value:   DefaultBranch,
		},
	},
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

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	args := rpc.InitArgs{
		Name:   c.Args().Get(0),
		Branch: c.String("branch"),
	}

	var reply rpc.InitReply
	if err := client.Call("Service.Init", &args, &reply); err != nil {
		return err
	}

	config := NewConfig(cwd)
	config.Name = args.Name
	config.Branch = args.Branch
	return config.Save()
}
