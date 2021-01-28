package main

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/rpc"
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

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	args := rpc.BranchArgs{
		Name:   config.Name,
		Branch: c.Args().Get(0),
		Head:   config.Index,
	}

	var action string
	switch {
	case c.Bool("create"):
		action = "Service.CreateBranch"
	case c.Bool("delete"):
		action = "Service.DeleteBranch"
	default:
		action = "Service.ListBranches"
	}

	var reply rpc.BranchReply
	if err := client.Call(action, &args, &reply); err != nil {
		return err
	}

	for branch := range reply.Branches {
		switch {
		case branch == config.Branch:
			fmt.Printf("* %s\n", branch)
		default:
			fmt.Println(branch)
		}
	}

	return nil
}
