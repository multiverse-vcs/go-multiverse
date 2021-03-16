package author

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/urfave/cli/v2"
)

// NewExportCommand returns a new command.
func NewExportCommand() *cli.Command {
	return &cli.Command{
		Name:      "export",
		Usage:     "Export an author",
		ArgsUsage: "[peer]",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowAppHelpAndExit(c, -1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := author.ExportArgs{
				Peer: c.Args().Get(0),
			}

			var reply author.ExportReply
			if err := client.Call("Author.Export", &args, &reply); err != nil {
				return err
			}

			fmt.Println(reply.Data)
			return nil
		},
	}
}
