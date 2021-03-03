package author

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/urfave/cli/v2"
)

// NewSelfCommand returns a new command.
func NewSelfCommand() *cli.Command {
	return &cli.Command{
		Name:  "self",
		Usage: "View your profile",
		Action: func(c *cli.Context) error {
			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := author.SelfArgs{}

			var reply author.SelfReply
			if err := client.Call("Author.Self", &args, &reply); err != nil {
				return err
			}

			fmt.Println(reply.PeerID.Pretty())
			return nil
		},
	}
}
