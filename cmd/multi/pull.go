package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var pullCommand = &cli.Command{
	Action:    pullAction,
	Name:      "pull",
	Usage:     "Merge changes into the current branch",
	ArgsUsage: "<cid>",
}

func pullAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
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

	if cfg.Head() != cfg.Index {
		return errors.New("index is behind head")
	}

	key, err := cfg.Key()
	if err != nil {
		return err
	}

	if err := node.Online(c.Context, key); err != nil {
		return err
	}

	fmt.Printf("fetching commit graph...\n")
	if err := merkledag.FetchGraph(c.Context, id, node.Dag); err != nil {
		return err
	}

	base, err := core.MergeBase(c.Context, node.Dag, cfg.Head(), id)
	if err != nil {
		return err
	}

	if base == id {
		return errors.New("local is ahead of remote")
	}

	merge, err := core.Merge(c.Context, repo, node.Dag, cfg.Head(), base, id)
	if err != nil {
		return err
	}

	return core.Write(c.Context, repo, node.Dag, "", merge)
}
