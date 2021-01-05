package main

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var statusCommand = &cli.Command{
	Action: statusAction,
	Name:   "status",
	Usage:  "Print changes",
}

func statusAction(c *cli.Context) error {
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

	args := rpc.StatusArgs{
		Root: config.Root,
		Head: head,
	}

	var reply rpc.StatusReply
	if err = client.Call("Service.Status", &args, &reply); err != nil {
		return err
	}

	fmt.Printf("Tracking changes on branch %s:\n", config.Branch)
	fmt.Printf("  (all files are automatically considered for commit)\n")
	fmt.Printf("  (to stop tracking files add rules to '%s')\n", IgnoreFile)

	for _, diff := range reply.Diffs {
		fmt.Println(diff)
	}

	return nil
}
