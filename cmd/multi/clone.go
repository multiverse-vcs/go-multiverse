package main

import (
	"os"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var cloneCommand = &cli.Command{
	Action:    cloneAction,
	Name:      "clone",
	Usage:     "Copy a repo",
	ArgsUsage: "<cid>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Directory name",
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch name",
		},
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"l"},
			Value:   -1,
			Usage:   "Fetch limit",
		},
	},
}

func cloneAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	args := rpc.CloneArgs{
		Cwd:    cwd,
		ID:     id,
		Limit:  c.Int("limit"),
		Name:   c.String("name"),
		Branch: c.String("branch"),
	}

	var reply rpc.CloneReply
	if err = client.Call("Service.Clone", &args, &reply); err != nil {
		return err
	}

	config := NewConfig(reply.Root, reply.Name)
	config.Branch = reply.Branch
	config.Branches = reply.Branches
	return config.Save()
}
