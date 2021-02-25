package branch

import (
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "branch",
		Usage: "Manage repository branches",
		Subcommands: []*cli.Command{
			NewListCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
		},
	}
}
