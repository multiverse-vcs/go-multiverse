package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var checkoutCommand = &cli.Command{
	Action:    checkoutAction,
	Name:      "checkout",
	Usage:     "Checkout committed files",
	ArgsUsage: "<cid>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "Force checkout with pending changes",
		},
	},
}

func checkoutAction(c *cli.Context) error {
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

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	changes, err := core.Status(c.Context, repo, node.Dag, cfg.Head())
	if err != nil {
		return err
	}

	if len(changes) > 0 && !c.Bool("force") {
		return errors.New("use the force flag to checkout with pending changes")
	}

	if err := core.Checkout(c.Context, repo, node.Dag, id); err != nil {
		return err
	}

	cfg.Index = id
	cfg.SetHead(id)

	return WriteConfig(root, &cfg)
}
