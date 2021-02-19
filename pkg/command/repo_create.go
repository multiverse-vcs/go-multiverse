package command

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/urfave/cli/v2"
)

// NewRepoCreateCommand returns a new command.
func NewRepoCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new repository",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			client, err := remote.NewClient(home)
			if err != nil {
				return ErrDialRPC
			}

			args := remote.CreateArgs{
				Name: c.Args().Get(0),
			}

			var reply remote.CreateReply
			if err := client.Call("Remote.Create", &args, &reply); err != nil {
				return err
			}

			fmt.Println(reply.Remote)
			return nil
		},
	}
}
