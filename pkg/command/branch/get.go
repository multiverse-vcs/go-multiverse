package branch

import (
	"errors"
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewGetCommand returns a new command.
func NewGetCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get branch settings",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
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
				fmt.Println(branch.Remote)
			default:
				return errors.New("invalid setting")
			}

			return nil
		},
	}
}
