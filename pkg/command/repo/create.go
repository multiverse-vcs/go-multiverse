package repo

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
	"github.com/urfave/cli/v2"
)

// NewCreateCommand returns a new command.
func NewCreateCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create a new repository",
		ArgsUsage: "[peer] [name]",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := repo.CreateArgs{
				Peer: c.Args().Get(0),
				Name: c.Args().Get(1),
			}

			var reply repo.CreateReply
			if err := client.Call("Repo.Create", &args, &reply); err != nil {
				return err
			}

			fmt.Println(reply.Remote)
			return nil
		},
	}
}
