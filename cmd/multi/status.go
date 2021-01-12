package main

import (
	"fmt"
	"os"
	"sort"

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

	ignore, err := config.Ignore()
	if err != nil {
		return err
	}

	args := rpc.StatusArgs{
		Root:   config.Root,
		Head:   config.Index,
		Ignore: ignore,
	}

	var reply rpc.StatusReply
	if err = client.Call("Service.Status", &args, &reply); err != nil {
		return err
	}

	paths := make([]string, 0)
	for path := range reply.Diffs {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	fmt.Printf("Tracking changes on branch %s:\n", config.Branch)
	fmt.Printf("  (all files are automatically considered for commit)\n")
	fmt.Printf("  (to stop tracking files add rules to '%s')\n", IgnoreFile)

	for _, p := range paths {
		switch reply.Diffs[p] {
		case rpc.StatusAdd:
			fmt.Printf("\tnew file: %s\n", p)
		case rpc.StatusRemove:
			fmt.Printf("\tdeleted:  %s\n", p)
		case rpc.StatusMod:
			fmt.Printf("\tmodified: %s\n", p)
		}
	}

	return nil
}
