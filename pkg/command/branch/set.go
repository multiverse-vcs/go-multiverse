package branch

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewSetCommand returns a new command.
func NewSetCommand() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "Set branch settings",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			cc, err := context.New(cwd)
			if err != nil {
				return err
			}

			branch := cc.Config.Branches[cc.Config.Branch]
			switch c.Args().Get(0) {
			case "remote":
				branch.Remote = c.Args().Get(1)
			default:
				return errors.New("invalid setting")
			}

			return cc.Config.Write()
		},
	}
}
