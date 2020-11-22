package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// NewCommitCommand returns a new commit command.
func NewCommitCommand() *cli.Command {
	return &cli.Command{
		Name:   "commit",
		Usage:  "record repo changes",
		Before: BeforeLoadContext,
		Action: func(c *cli.Context) error {
			id, err := cmdctx.Commit("")
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("%s%s%s\n", ColorYellow, id.String(), ColorReset)
			return nil
		},
	}
}
