package cmd

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/urfave/cli/v2"
)

// LogDateFormat is the format used when logging the commit.
const LogDateFormat = "Mon Jan 2 15:04:05 2006 -0700"

// NewLogCommand returns a new log command.
func NewLogCommand() *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "print change history",
		Action: func(c *cli.Context) error {
			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if !cfg.Head.Defined() {
				return nil
			}

			cb := func(id cid.Cid, commit *object.Commit) bool {
				fmt.Printf("%s%s", ColorYellow, id.String())

				if id == cfg.Head {
					fmt.Printf(" (%sHEAD%s)", ColorRed, ColorYellow)
				}

				if id == cfg.Base {
					fmt.Printf(" (%sBASE%s)", ColorGreen, ColorYellow)
				}

				fmt.Printf("%s\n", ColorReset)
				fmt.Printf("Date: %s\n", commit.Date.Format(LogDateFormat))
				fmt.Printf("\n\t%s\n\n", commit.Message)
				return true
			}

			if _, err := core.Walk(c.Context, store, cfg.Head, cb); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
