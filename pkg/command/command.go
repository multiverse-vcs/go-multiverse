package command

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/repo"
	"github.com/urfave/cli/v2"
)

// NewApp returns a new cli app.
func NewApp() *cli.App {
	return &cli.App{
		Name:        "multi",
		HelpName:    "multi",
		Usage:       "Multiverse command line interface",
		Description: `Multiverse is a decentralized version control system for peer-to-peer software development.`,
		Version:     "0.0.5",
		Authors: []*cli.Author{
			{Name: "Keenan Nemetz", Email: "keenan.nemetz@pm.me"},
		},
		Commands: []*cli.Command{
			NewInitCommand(),
			NewCommitCommand(),
			NewRemoteCommand(),
			NewPushCommand(),
			NewStatusCommand(),
			NewDaemonCommand(),
			repo.NewCommand(),
		},
	}
}

// Execute runs the cli app.
func Execute() {
	if err := NewApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
