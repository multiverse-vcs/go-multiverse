package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

// DotDir is the name of the root directory.
const DotDir = ".multiverse"

var fs = afero.NewOsFs()

var app = &cli.App{
	Name:     "multi",
	HelpName: "multi",
	Usage:    "Decentralized Version Control System",
	Description: `Multiverse is a decentralized version control system
that enables peer-to-peer software development.`,
	Version: "0.0.1",
	Authors: []*cli.Author{
		{Name: "Keenan Nemetz", Email: "keenan.nemetz@pm.me"},
	},
	Commands: []*cli.Command{
		branchCommand,
		checkoutCommand,
		commitCommand,
		historyCommand,
		initCommand,
		pullCommand,
		pushCommand,
		statusCommand,
		switchCommand,
	},
}

func init() {
	core.IgnoreRules = append(core.IgnoreRules, DotDir)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Root finds the repo root by searching parent directories.
func Root(root string) (string, error) {
	path := filepath.Join(root, DotDir)

	info, err := fs.Stat(path)
	if err == nil && info.IsDir() {
		return root, nil
	}

	parent := filepath.Dir(root)
	if parent == root {
		return "", errors.New("repo not found")
	}

	return Root(parent)
}
