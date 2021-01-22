package main

import (
	"os"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var mergeCommand = &cli.Command{
	Action:    mergeAction,
	Name:      "merge",
	Usage:     "Merge commits",
	ArgsUsage: "<cid>",
}

func mergeAction(c *cli.Context) error {
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

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	args := rpc.MergeArgs{
		Repo:   config.Repo,
		Branch: config.Branch,
		Root:   config.Root,
		Index:  config.Index,
		ID:     id,
	}

	var reply rpc.MergeReply
	if err = client.Call("Service.Merge", &args, &reply); err != nil {
		return err
	}

	config.Repo = reply.Repo
	config.Index = reply.Index
	return config.Save()
}
