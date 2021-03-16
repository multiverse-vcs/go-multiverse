package repo

import (
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
	"github.com/urfave/cli/v2"
)

// NewDeleteCommand returns a new command.
func NewDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete an existing repository",
		ArgsUsage: "[peer] [name]",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := repo.DeleteArgs{
				Peer: c.Args().Get(0),
				Name: c.Args().Get(1),
			}

			var reply repo.DeleteReply
			if err := client.Call("Repo.Delete", &args, &reply); err != nil {
				return err
			}

			return nil
		},
	}
}
