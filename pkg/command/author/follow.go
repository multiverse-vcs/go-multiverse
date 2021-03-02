package author

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/urfave/cli/v2"
)

// NewFollowCommand returns a new command.
func NewFollowCommand() *cli.Command {
	return &cli.Command{
		Name:  "follow",
		Usage: "Subscribe to author updates",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			peerID, err := peer.Decode(c.Args().Get(0))
			if err != nil {
				return err
			}

			args := author.FollowArgs{
				PeerID: peerID,
			}

			var reply author.FollowReply
			if err := client.Call("Author.Follow", &args, &reply); err != nil {
				return err
			}

			return nil
		},
	}
}
