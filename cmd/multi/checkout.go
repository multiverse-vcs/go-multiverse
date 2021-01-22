package main

import (
	"os"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var checkoutCommand = &cli.Command{
	Action:    checkoutAction,
	Name:      "checkout",
	Usage:     "Checkout committed files",
	ArgsUsage: "<cid>",
}

func checkoutAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := LoadConfig(cwd)
	if err != nil {
		return err
	}

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	id, err := cid.Parse(c.Args().Get(0))
	if err != nil {
		return err
	}

	args := rpc.CheckoutArgs{
		Repo:   config.Repo,
		Branch: config.Branch,
		Root:   config.Root,
		Index:  config.Index,
		ID:     id,
	}

	var reply rpc.CheckoutReply
	if err = client.Call("Service.Checkout", &args, &reply); err != nil {
		return err
	}

	config.Index = id
	return config.Save()
}
