package main

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
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
	store, err := openStore()
	if err != nil {
		return err
	}

	cfg, err := store.ReadConfig()
	if err != nil {
		return err
	}

	var parents []cid.Cid
	if cfg.Head().Defined() {
		parents = append(parents, cfg.Head())
	}

	id, err := core.Commit(c.Context, store, c.String("message"), parents...)
	if err != nil {
		return err
	}

	cfg.Index = id
	cfg.SetHead(id)

	if err := store.WriteConfig(cfg); err != nil {
		return err
	}

	fmt.Println(id.String())
	return nil
}
