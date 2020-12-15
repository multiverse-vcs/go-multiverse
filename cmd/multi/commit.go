package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var commitCommand = &cli.Command{
	Action: commitAction,
	Name:   "commit",
	Usage:  "Record repo changes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
			Value:   "",
			Usage:   "Description of changes",
		},
	},
}

func commitAction(c *cli.Context) error {
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

	var parents []cid.Cid
	if cfg.Head().Defined() {
		parents = append(parents, cfg.Head())
	}

	id, err := core.Commit(c.Context, repo, node.Dag, c.String("message"), parents...)
	if err != nil {
		return err
	}
	fmt.Println(id.String())

	cfg.Index = id
	cfg.SetHead(id)

	return WriteConfig(root, &cfg)
}
