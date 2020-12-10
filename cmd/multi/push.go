package main

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/remote"
	"github.com/urfave/cli/v2"
)

var pushCommand = &cli.Command{
	Action: pushAction,
	Name:   "push",
	Usage:  "Upload changes to a remote",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "remote",
			Aliases: []string{"r"},
			Value:   config.DefaultRemote,
			Usage:   "Remote to push changes to",
		},
	},
}

func pushAction(c *cli.Context) error {
	store, err := openStore()
	if err != nil {
		return err
	}

	cfg, err := store.ReadConfig()
	if err != nil {
		return err
	}

	url, ok := cfg.Remotes[c.String("remote")]
	if !ok {
		return cli.Exit("remote does not exist", 1)
	}

	if !cfg.Head().Defined() {
		return cli.Exit("nothing to push", 1)
	}

	client := remote.NewRemote(url)
	if err := client.Upload(c.Context, store, cfg.Head()); err != nil {
		return err
	}

	fmt.Println("changes pushed successfully!")
	return nil
}
