package main

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var commitCommand = &cli.Command{
	Action: commitAction,
	Name:   "commit",
	Usage:  "Record changes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
			Usage:   "Description",
		},
	},
}

func commitAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := LoadConfig(cwd)
	if err != nil {
		return err
	}

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	ignore, err := config.Ignore()
	if err != nil {
		return err
	}

	args := rpc.CommitArgs{
		Root:    config.Root,
		Ignore:  ignore,
		Repo:    config.Repo,
		Branch:  config.Branch,
		Parent:  config.Index,
		Message: c.String("message"),
	}

	var reply rpc.CommitReply
	if err = client.Call("Service.Commit", &args, &reply); err != nil {
		return err
	}

	config.Repo = reply.Repo
	config.Index = reply.Index
	return config.Save()
}
