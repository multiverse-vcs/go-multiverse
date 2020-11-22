// Package command implements the Multiverse CLI.
package cmd

import (
	"os"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/urfave/cli/v2"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
)

// cmdctx is the command context for all commands
var cmdctx *Context

// NewApp returns a new cli app.
func NewApp() *cli.App {
	return &cli.App{
		Name:     "multi",
		HelpName: "multi",
		Usage:    "decentralized version control system",
		Description: `Multiverse is a decentralized version control system
that enables peer-to-peer software development.`,
		Version: "0.0.1",
		Authors: []*cli.Author{
			{Name: "Keenan Nemetz", Email: "keenan.nemetz@pm.me"},
		},
		Commands: []*cli.Command{
			NewCommitCommand(),
			NewInitCommand(),
			NewLogCommand(),
		},
	}
}

// BeforeLoadContext is used as a before func to load cmdctx.
func BeforeLoadContext(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	cmdctx, err = LoadContext(osfs.New(cwd), c.Context)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}

// AfterWriteConfig is used as an after func to write config.
func AfterWriteConfig(c *cli.Context) error {
	return nil
}
