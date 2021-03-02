package remote

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewCreateCommand returns a new command.
func NewCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new remote",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
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
			if _, ok := cc.Config.Remotes[name]; ok {
				return errors.New("remote already exists")
			}

			cc.Config.Remotes[name] = c.Args().Get(1)
			return cc.Config.Write()
		},
	}
}
