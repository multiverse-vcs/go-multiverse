package main

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var tagCommand = &cli.Command{
	Action:    tagAction,
	Name:      "tag",
	Usage:     "List, create, or delete tags",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete tag",
		},
		&cli.BoolFlag{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create tag",
		},
	},
}

func tagAction(c *cli.Context) error {
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

	args := rpc.TagArgs{
		Name: config.Name,
		Tag:  c.Args().Get(0),
		Head: config.Index,
	}

	var action string
	switch {
	case c.Bool("create"):
		action = "Service.CreateTag"
	case c.Bool("delete"):
		action = "Service.DeleteTag"
	default:
		action = "Service.ListTags"
	}

	var reply rpc.TagReply
	if err = client.Call(action, &args, &reply); err != nil {
		return err
	}

	for tag := range reply.Tags {
		fmt.Println(tag)
	}

	return nil
}
