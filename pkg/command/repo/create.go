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
		Name:  "create",
		Usage: "Create a new repository",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := repo.CreateArgs{
				Name: c.Args().Get(0),
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
