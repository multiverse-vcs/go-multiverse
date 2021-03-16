package remote

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewRemoveCommand returns a new command.
func NewRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "Remove an existing remote",
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

			alias := c.Args().Get(0)
			if _, ok := cc.Config.Remotes[alias]; !ok {
				return errors.New("remote does not exists")
			}

			delete(cc.Config.Remotes, alias)
			return cc.Config.Write()
		},
	}
}
