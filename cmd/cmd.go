package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/storage"
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

// NewApp returns a new cli app.
func NewApp() *cli.App {
	return &cli.App{
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
			NewBranchCommand(),
			NewCheckoutCommand(),
			NewCloneCommand(),
			NewCommitCommand(),
			NewInitCommand(),
			NewLogCommand(),
			NewStatusCommand(),
			NewSwapCommand(),
			NewSwitchCommand(),
		},
	}
}

// Root finds the repo root by searching parent directories.
func Root(root string) (string, error) {
	path := filepath.Join(root, storage.DotDir)

	info, err := os.Lstat(path)
	if err == nil && info.IsDir() {
		return root, nil
	}

	parent := filepath.Dir(root)
	if parent == root {
		return "", errors.New("repo not found")
	}

	return Root(parent)
}

// Store returns the repo store from the current directory.
func Store() (*storage.Store, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root, err := Root(cwd)
	if err != nil {
		return nil, err
	}

	return storage.NewOsStore(root)
}
