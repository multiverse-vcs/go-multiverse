package branch

import (
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "branch",
		Usage: "List, create, or delete branches",
		Subcommands: []*cli.Command{
			NewListCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
			NewSetCommand(),
			NewGetCommand(),
		},
	}
}
