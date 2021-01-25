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
			Name:    "dir",
			Aliases: []string{"d"},
			Usage:   "Directory name",
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch name",
			Value:   "default",
		},
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"l"},
			Usage:   "Fetch limit",
			Value:   -1,
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
		Dir:    c.String("dir"),
		ID:     id,
		Limit:  c.Int("limit"),
		Branch: c.String("branch"),
	}

	var reply rpc.CloneReply
	if err = client.Call("Service.Clone", &args, &reply); err != nil {
		return err
	}

	config := NewConfig(reply.Root)
	config.Branch = c.String("branch")
	config.Index = reply.ID
	return config.Save()
}
