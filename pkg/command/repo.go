package command

import (
	"github.com/urfave/cli/v2"
)

// NewRepoCommand returns a new command.
func NewRepoCommand() *cli.Command {
	return &cli.Command{
		Name:  "repo",
		Usage: "Manage remote repositories",
		Subcommands: []*cli.Command{
			NewRepoCreateCommand(),
		},
	}
}
