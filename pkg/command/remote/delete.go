package remote

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewDeleteCommand returns a new command.
func NewDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Create an existing remote",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowAppHelpAndExit(c, -1)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			cc, err := context.New(cwd)
			if err != nil {
				return err
			}

			name := c.Args().Get(0)
			if _, ok := cc.Config.Remotes[name]; !ok {
				return errors.New("remote does not exists")
			}

			delete(cc.Config.Remotes, name)
			return cc.Config.Write()
		},
	}
}
