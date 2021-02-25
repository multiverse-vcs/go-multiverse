package branch

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewListCommand returns a new command.
func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all branches",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			ctx, err := context.New(cwd)
			if err != nil {
				return err
			}

			for branch := range ctx.Config.Repository.Branches {
				if branch == ctx.Config.Branch {
					fmt.Print("* ")
				}

				fmt.Println(branch)
			}

			return nil
		},
	}
}
