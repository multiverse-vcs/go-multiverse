package author

import (
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "author",
		Usage: "Manage author profiles",
		Subcommands: []*cli.Command{
			NewSelfCommand(),
			NewListCommand(),
			NewViewCommand(),
			NewFollowCommand(),
			NewUnfollowCommand(),
		},
	}
}
