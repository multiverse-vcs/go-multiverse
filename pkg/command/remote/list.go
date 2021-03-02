package remote

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
		Usage: "List all remotes",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			cc, err := context.New(cwd)
			if err != nil {
				return err
			}

			for name, path := range cc.Config.Remotes {
				fmt.Printf("%-32s%s\n", name, path)
			}

			return nil
		},
	}
}
