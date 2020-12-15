package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/urfave/cli/v2"
)

// DateFormat is the format used when logging the commit.
const DateFormat = "Mon Jan 2 15:04:05 2006 -0700"

var historyCommand = &cli.Command{
	Action:  historyAction,
	Name:    "history",
	Aliases: []string{"log"},
	Usage:   "Print change history",
}

func historyAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := Root(cwd)
	if err != nil {
		return err
	}

	root := filepath.Join(path, DotDir)

	node, err := node.NewNode(root)
	if err != nil {
		return err
	}

	var cfg Config
	if err := ReadConfig(root, &cfg); err != nil {
		return err
	}

	if !cfg.Head().Defined() {
		return nil
	}

	cb := func(id cid.Cid, commit *object.Commit) bool {
		fmt.Printf("%s", id.String())

		if id == cfg.Head() {
			fmt.Printf(" (HEAD)")
		}

		if id == cfg.Index {
			fmt.Printf(" (INDEX)")
		}

		fmt.Printf("\nDate: %s\n\n", commit.Date.Format(DateFormat))

		if len(commit.Message) > 0 {
			fmt.Printf("\t%s\n\n", commit.Message)
		}

		return true
	}

	if _, err := core.Walk(c.Context, node.Dag, cfg.Head(), cb); err != nil {
		return err
	}

	return nil
}
