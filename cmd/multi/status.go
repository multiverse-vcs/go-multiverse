package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var statusCommand = &cli.Command{
	Action: statusAction,
	Name:   "status",
	Usage:  "Print repo status",
}

func statusAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := Root(cwd)
	if err != nil {
		return err
	}

	repo := afero.NewBasePathFs(fs, path)
	root := filepath.Join(path, DotDir)

	node, err := node.NewNode(root)
	if err != nil {
		return err
	}

	var cfg Config
	if err := ReadConfig(root, &cfg); err != nil {
		return err
	}

	changes, err := core.Status(c.Context, repo, node.Dag, cfg.Head())
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	fmt.Printf("Tracking changes on branch %s:\n", cfg.Branch)
	fmt.Printf("  (all files are automatically considered for commit)\n")
	fmt.Printf("  (to stop tracking files add rules to '%s')\n", core.IgnoreFile)

	set := make(map[string]dagutils.ChangeType)
	for _, change := range changes {
		if _, ok := set[change.Path]; ok {
			set[change.Path] = dagutils.Mod
		} else {
			set[change.Path] = change.Type
		}
	}

	paths := make([]string, 0)
	for path := range set {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, p := range paths {
		switch set[p] {
		case dagutils.Add:
			fmt.Printf("\tnew file: %s\n", p)
		case dagutils.Remove:
			fmt.Printf("\tdeleted:  %s\n", p)
		case dagutils.Mod:
			fmt.Printf("\tmodified: %s\n", p)
		}
	}

	return nil
}
