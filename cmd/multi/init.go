package main

import (
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/urfave/cli/v2"
)

var initCommand = &cli.Command{
	Action: initAction,
	Name:   "init",
	Usage:  "Initialize a new repo",
}

func initAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(cwd, storage.DotDir)
	if err := fs.Mkdir(path, 0755); err != nil {
		return err
	}

	store, err := openStore()
	if err != nil {
		return err
	}

	return store.Initialize()
}
