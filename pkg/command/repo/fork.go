package repo

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
	"github.com/urfave/cli/v2"
)

// NewForkCommand returns a new command.
func NewForkCommand() *cli.Command {
	return &cli.Command{
		Name:  "fork",
		Usage: "Copy an existing repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "New repository name",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := repo.ForkArgs{
				Name: c.String("name"),
				Remote: c.Args().Get(0),
			}

			var reply repo.ForkReply
			if err := client.Call("Repo.Fork", &args, &reply); err != nil {
				return err
			}

			fmt.Println(reply.Remote)
			return nil
		},
	}
}
