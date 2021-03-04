package repo

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
)

// NewImportCommand returns a new command.
func NewImportCommand() *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: "Import an external repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Usage: "Repository url",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Repository path",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			if !c.IsSet("url") && !c.IsSet("path") {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := repo.ImportArgs{
				Name: c.Args().Get(0),
				URL:  c.String("url"),
				Path: c.String("path"),
			}

			var reply repo.ImportReply
			if err := client.Call("Repo.Import", &args, &reply); err != nil {
				return err
			}

			fmt.Println(reply.Remote)
			return nil
		},
	}
}
