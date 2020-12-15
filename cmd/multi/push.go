package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/node"
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
			Value:   "local",
			Usage:   "Remote to push changes to",
		},
	},
}

func pushAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := Root(cwd)
	if err != nil {
		return err
	}

	root := filepath.Join(path, DotDir)

	node, err := node.NewNode(root)
	if err != nil {
		return err
	}

	var cfg Config
	if err := ReadConfig(root, &cfg); err != nil {
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
	if err := client.Upload(c.Context, node.Dag, cfg.Head()); err != nil {
		return err
	}

	fmt.Println("changes pushed successfully!")
	return nil
}
