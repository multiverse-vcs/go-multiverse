package remote

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewAddCommand returns a new command.
func NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new remote",
		ArgsUsage: "[alias] [peer] [name]",
		Action: func(c *cli.Context) error {
			if c.NArg() != 3 {
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

			alias := c.Args().Get(0)
			if _, ok := cc.Config.Remotes[alias]; ok {
				return errors.New("remote already exists")
			}

			cc.Config.Remotes[alias] = &context.Remote{
				Peer: c.Args().Get(1),
				Name: c.Args().Get(2),
			}

			return cc.Config.Write()
		},
	}
}
