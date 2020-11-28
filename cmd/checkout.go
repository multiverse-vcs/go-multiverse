package cmd

import (
	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/urfave/cli/v2"
)

// NewCheckoutCommand returns a new command.
func NewCheckoutCommand() *cli.Command {
	return &cli.Command{
		Name:      "checkout",
		Usage:     "checkout committed files",
		ArgsUsage: "<commit-cid>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Value:   false,
				Usage:   "force checkout with pending changes",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			id, err := cid.Parse(c.Args().Get(0))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			changes, err := core.Status(c.Context, store, cfg.Head())
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if len(changes) > 0 && !c.Bool("force") {
				return cli.Exit("use the force flag to checkout with pending changes", 1)
			}

			if err := core.Checkout(c.Context, store, id); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg.Index = id
			cfg.SetHead(id)

			if err := store.WriteConfig(cfg); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
