package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// NewRemoteCommand returns a new cli command.
func NewRemoteCommand() *cli.Command {
	return &cli.Command{
		Name:  "remote",
		Usage: "Get or set the repository remote",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			ctx, err := NewContext(cwd)
			if err != nil {
				return err
			}

			// TODO add some verification of remote

			switch c.NArg() {
			case 0:
				fmt.Println(ctx.Config.Remote)
				return nil
			case 1:
				ctx.Config.Remote = c.Args().Get(0)
				return ctx.Config.Write()
			default:
				cli.ShowAppHelpAndExit(c, -1)
				return nil
			}
		},
	}
}
