package main

import (
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
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

	store, err := openStore()
	if err != nil {
		return err
	}

	cfg, err := store.ReadConfig()
	if err != nil {
		return err
	}

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	changes, err := core.Status(c.Context, store, cfg.Head())
	if err != nil {
		return err
	}

	if len(changes) > 0 && !c.Bool("force") {
		return errors.New("use the force flag to checkout with pending changes")
	}

	if err := core.Checkout(c.Context, store, id); err != nil {
		return err
	}

	cfg.Index = id
	cfg.SetHead(id)

	return store.WriteConfig(cfg)
}
