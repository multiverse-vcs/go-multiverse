package main

import (
	"os"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var pullCommand = &cli.Command{
	Action:    pullAction,
	Name:      "pull",
	Usage:     "Merge changes",
	ArgsUsage: "<cid>",
}

func pullAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

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

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	args := rpc.PullArgs{
		Root: config.Root,
		Head: head,
		ID:   id,
	}

	var reply rpc.PullReply
	if err = client.Call("Service.Pull", &args, &reply); err != nil {
		return err
	}

	config.SetHead(reply.ID)
	return config.Save()
}
