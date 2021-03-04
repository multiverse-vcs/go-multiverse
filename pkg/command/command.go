package command

import (
	"time"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/author"
	"github.com/multiverse-vcs/go-multiverse/pkg/command/branch"
	"github.com/multiverse-vcs/go-multiverse/pkg/command/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/command/repo"
	"github.com/urfave/cli/v2"
)

// NewApp returns a new cli app.
func NewApp() *cli.App {
	return &cli.App{
		Compiled:    time.Now(),
		Name:        "multi",
		HelpName:    "multi",
		Usage:       "Multiverse command line interface",
		Description: `Multiverse is a decentralized version control system for peer-to-peer software development.`,
		Authors: []*cli.Author{
			{Name: "Keenan Nemetz", Email: "keenan.nemetz@pm.me"},
		},
		Commands: []*cli.Command{
			NewInitCommand(),
			NewCommitCommand(),
			NewCheckoutCommand(),
			NewSwitchCommand(),
			NewPushCommand(),
			NewPullCommand(),
			NewStatusCommand(),
			NewLogCommand(),
			branch.NewCommand(),
			remote.NewCommand(),
			repo.NewCommand(),
			author.NewCommand(),
			NewDaemonCommand(),
		},
	}
}
