package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var switchCommand = &cli.Command{
	Action:    switchAction,
	Name:      "switch",
	Usage:     "Change branches",
	ArgsUsage: "<name>",
}

func switchAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

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

	name := c.Args().Get(0)
	if cfg.Branch == name {
		return errors.New("already on branch")
	}

	branch, err := cfg.GetBranch(name)
	if err != nil {
		return err
	}

	fmt.Println("stashing changes...")
	stash, err := core.Worktree(c.Context, repo, node.Dag)
	if err != nil {
		return err
	}

	var id cid.Cid = branch.Head
	if branch.Stash.Defined() {
		id = branch.Stash
	}

	fmt.Println("checking out branch...")
	if err := core.Checkout(c.Context, repo, node.Dag, id); err != nil {
		return err
	}

	cfg.SetStash(stash.Cid())
	cfg.Branch = name
	cfg.Index = id

	return WriteConfig(root, &cfg)
}
