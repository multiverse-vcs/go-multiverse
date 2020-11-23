package cmd

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/urfave/cli/v2"
)

// LogDateFormat is the format used when logging the commit.
const LogDateFormat = "Mon Jan 2 15:04:05 2006 -0700"

// NewLogCommand returns a new log command.
func NewLogCommand() *cli.Command {
	return &cli.Command{
		Name:   "log",
		Usage:  "print change history",
		Before: BeforeLoadContext,
		Action: func(c *cli.Context) error {
			cb := func(id cid.Cid, commit *object.Commit) bool {
				fmt.Printf("%s%s", ColorYellow, id.String())

				if id == cmdctx.Config.Head {
					fmt.Printf(" (%sHEAD%s)", ColorRed, ColorYellow)
				}

				if id == cmdctx.Config.Base {
					fmt.Printf(" (%sBASE%s)", ColorGreen, ColorYellow)
				}

				fmt.Printf("%s\n", ColorReset)
				fmt.Printf("Date: %s\n", commit.Date.Format(LogDateFormat))
				fmt.Printf("\n\t%s\n\n", commit.Message)
				return true
			}

			if !cmdctx.Config.Head.Defined() {
				return nil
			}

			_, err := cmdctx.Walk(cmdctx.Config.Head, cb)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
