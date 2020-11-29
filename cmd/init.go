package cmd

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/urfave/cli/v2"
)

// NewInitCommand returns a new init command.
func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "create a new repo",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if _, err := storage.InitOsStore(cwd); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
