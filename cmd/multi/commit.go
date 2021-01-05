package main

import (
	"os"

	"github.com/ipfs/go-cid"
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

	head, err := config.Head()
	if err != nil {
		return err
	}

	var parents []cid.Cid
	if head.Defined() {
		parents = append(parents, head)
	}

	ignore, err := config.Ignore()
	if err != nil {
		return err
	}

	args := rpc.CommitArgs{
		Root:    config.Root,
		Ignore:  ignore,
		Name:    config.Name,
		Branch:  config.Branch,
		Parents: parents,
		Message: c.String("message"),
	}

	var reply rpc.CommitReply
	if err = client.Call("Service.Commit", &args, &reply); err != nil {
		return err
	}

	config.SetHead(reply.ID)
	return config.Save()
}
