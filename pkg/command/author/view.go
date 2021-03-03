package author

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/urfave/cli/v2"
)

// NewViewCommand returns a new command.
func NewViewCommand() *cli.Command {
	return &cli.Command{
		Name:  "view",
		Usage: "View profile details",
		Action: func(c *cli.Context) error {
			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			peerID, err := peer.Decode(c.Args().Get(0))
			if err != nil {
				return err
			}

			args := author.SearchArgs{
				PeerID: peerID,
			}

			var reply author.SearchReply
			if err := client.Call("Author.Self", &args, &reply); err != nil {
				return err
			}

			fmt.Println("")
			fmt.Printf("Name%28sCID\n", "")
			fmt.Printf("----%28s---\n", "")

			for name, id := range reply.Author.Repositories {
				fmt.Printf("%-32s%s\n", name, id.String())
			}

			return nil
		},
	}
}
