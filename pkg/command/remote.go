package command

import (
	"os"

	"github.com/urfave/cli/v2"
)

// NewRemoteCommand returns a new cli command.
func NewRemoteCommand() *cli.Command {
	return &cli.Command{
		Name:  "remote",
		Usage: "Set the repository remote",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			repo, err := NewContext(cwd)
			if err != nil {
				return err
			}

			repo.Config.Remote = c.Args().Get(0)
			return repo.Config.Write()
		},
	}
}
