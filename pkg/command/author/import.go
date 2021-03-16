package author

import (
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/urfave/cli/v2"
)

// NewImportCommand returns a new command.
func NewImportCommand() *cli.Command {
	return &cli.Command{
		Name:      "import",
		Usage:     "Import an author",
		ArgsUsage: "[data]",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowAppHelpAndExit(c, -1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := author.ImportArgs{
				Data: c.Args().Get(0),
			}

			var reply author.ImportReply
			if err := client.Call("Author.Import", &args, &reply); err != nil {
				return err
			}

			return nil
		},
	}
}
