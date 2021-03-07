package repo

import (
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "repo",
		Usage: "Manage remote repositories",
		Subcommands: []*cli.Command{
			NewCreateCommand(),
			NewForkCommand(),
			NewListCommand(),
			NewDeleteCommand(),
			NewImportCommand(),
		},
	}
}
