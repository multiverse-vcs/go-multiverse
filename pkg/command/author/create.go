package author

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/urfave/cli/v2"
)

// NewCreateCommand returns a new command.
func NewCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new author",
		Action: func(c *cli.Context) error {
			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := author.CreateArgs{}

			var reply author.CreateReply
			if err := client.Call("Author.Create", &args, &reply); err != nil {
				return err
			}

			fmt.Println(reply.Peer)
			return nil
		},
	}
}
