package remote

import (
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "remote",
		Usage: "List, create, or delete remotes",
		Subcommands: []*cli.Command{
			NewListCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
		},
	}
}
