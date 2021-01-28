package main

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

const logDateFormat = "Mon Jan 02 15:04:05 2006 -0700"

var logCommand = &cli.Command{
	Action: logAction,
	Name:   "log",
	Usage:  "Print repo history",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"l"},
			Usage:   "Log limit",
			Value:   -1,
		},
	},
}

func logAction(c *cli.Context) error {
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

	args := rpc.LogArgs{
		Name:   config.Name,
		Branch: config.Branch,
		Limit:  c.Int("limit"),
	}

	var reply rpc.LogReply
	if err := client.Call("Service.Log", &args, &reply); err != nil {
		return err
	}

	for i, c := range reply.Commits {
		fmt.Printf("commit %s\n", reply.IDs[i].String())
		fmt.Printf("Date:  %s\n", c.Date.Format(logDateFormat))
		fmt.Printf("\n\t%s\n\n", c.Message)
	}

	return nil
}
