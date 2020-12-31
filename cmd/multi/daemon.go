package main

import (
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var daemonCommand = &cli.Command{
	Action: daemonAction,
	Name:   "daemon",
	Usage:  "Starts a client",
}

func daemonAction(c *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	root := filepath.Join(home, ".multiverse")
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	node, err := node.New(c.Context, root)
	if err != nil {
		return err
	}

	return rpc.ListenAndServe(node)
}
