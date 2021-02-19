package command

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
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

			path := remote.Path(c.Args().Get(0))
			if err = path.Verify(); err != nil {
				return err
			}

			repo.Config.Remote = path
			return repo.Config.Write()
		},
	}
}
