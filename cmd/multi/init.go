package main

import (
	"errors"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var initCommand = &cli.Command{
	Action:    initAction,
	Name:      "init",
	Usage:     "Create a repo",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "cid",
			Aliases: []string{"c"},
			Usage:   "CID",
		},
	},
}

func initAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := FindConfig(cwd); err == nil {
		return errors.New("repo already exists")
	}

	config := DefaultConfig(cwd, c.Args().Get(0))
	if !c.IsSet("cid") {
		return config.Save()
	}

	id, err := cid.Parse(c.String("cid"))
	if err != nil {
		return err
	}

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	args := rpc.CheckoutArgs{
		Root: config.Root,
		ID:   id,
	}

	var reply rpc.CheckoutReply
	if err = client.Call("Service.Checkout", &args, &reply); err != nil {
		return err
	}

	config.SetHead(id)
	return config.Save()
}
