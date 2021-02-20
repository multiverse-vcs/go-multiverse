package repo

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
	"github.com/urfave/cli/v2"
)

// NewListCommand returns a new command.
func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all repositories",
		Action: func(c *cli.Context) error {
			client, err := rpc.NewClient()
			if err != nil {
				return rpc.ErrDialRPC
			}

			args := repo.ListArgs{}

			var reply repo.ListReply
			if err := client.Call("Repo.List", &args, &reply); err != nil {
				return err
			}

			for name := range reply.Repositories {
				fmt.Println(name)
			}

			return nil
		},
	}
}
