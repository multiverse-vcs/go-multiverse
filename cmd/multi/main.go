package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
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
		cloneCommand,
		commitCommand,
		historyCommand,
		initCommand,
		pullCommand,
		pushCommand,
		statusCommand,
		switchCommand,
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// findRoot finds the repo root by searching parent directories.
func findRoot(root string) (string, error) {
	path := filepath.Join(root, storage.DotDir)

	info, err := fs.Stat(path)
	if err == nil && info.IsDir() {
		return root, nil
	}

	parent := filepath.Dir(root)
	if parent == root {
		return "", errors.New("repo not found")
	}

	return findRoot(parent)
}

// openStore returns the repo store from the current directory.
func openStore() (*storage.Store, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root, err := findRoot(cwd)
	if err != nil {
		return nil, err
	}

	return storage.NewStore(fs, root)
}
