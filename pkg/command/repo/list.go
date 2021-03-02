package repo

import (
	"fmt"

	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
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
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			args := author.SelfArgs{}

			var reply author.SelfReply
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
