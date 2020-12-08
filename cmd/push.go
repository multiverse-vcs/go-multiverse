package cmd

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/remote"
	"github.com/urfave/cli/v2"
)

// NewPushCommand returns a new status command.
func NewPushCommand() *cli.Command {
	return &cli.Command{
		Name:  "push",
		Usage: "Upload changes to a remote",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "remote",
				Aliases: []string{"r"},
				Value:   config.DefaultRemote,
				Usage:   "Remote to push changes to",
			},
		},
		Action: func(c *cli.Context) error {
			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
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
				return cli.Exit(err.Error(), 1)
			}

			fmt.Println("changes pushed successfully!")
			return nil
		},
	}
}
