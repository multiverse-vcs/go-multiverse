package branch

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
		Usage: "Create a new branch",
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
			if _, ok := cc.Config.Branches[name]; ok {
				return errors.New("branch already exists")
			}

			branch := cc.Config.Branches[cc.Config.Branch]

			cc.Config.Branches[name] = &context.Branch{
				Head:   branch.Head,
				Stash:  branch.Stash,
				Remote: branch.Remote,
			}

			return cc.Config.Write()
		},
	}
}
