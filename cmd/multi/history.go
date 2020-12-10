package main

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/urfave/cli/v2"
)

// DateFormat is the format used when logging the commit.
const DateFormat = "Mon Jan 2 15:04:05 2006 -0700"

var historyCommand = &cli.Command{
	Action: historyAction,
	Name:   "history",
	Usage:  "Print change history",
}

func historyAction(c *cli.Context) error {
	store, err := openStore()
	if err != nil {
		return err
	}

	cfg, err := store.ReadConfig()
	if err != nil {
		return err
	}

	if !cfg.Head().Defined() {
		return nil
	}

	cb := func(id cid.Cid, commit *object.Commit) bool {
		fmt.Printf("%s%s", ColorYellow, id.String())

		if id == cfg.Head() {
			fmt.Printf(" (%sHEAD%s)", ColorRed, ColorYellow)
		}

		if id == cfg.Index {
			fmt.Printf(" (%sINDEX%s)", ColorGreen, ColorYellow)
		}

		fmt.Printf("%s\n", ColorReset)
		fmt.Printf("Date: %s\n", commit.Date.Format(DateFormat))

		if len(commit.Message) > 0 {
			fmt.Printf("\n\t%s\n", commit.Message)
		}

		fmt.Printf("\n")
		return true
	}

	if _, err := core.Walk(c.Context, store, cfg.Head(), cb); err != nil {
		return err
	}

	return nil
}
