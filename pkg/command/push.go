package command

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
	"github.com/urfave/cli/v2"
)

// NewPushCommand returns a new cli command.
func NewPushCommand() *cli.Command {
	return &cli.Command{
		Name:  "push",
		Usage: "Update a remote repository",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			ctx, err := NewContext(cwd)
			if err != nil {
				return err
			}

			client, err := rpc.NewClient()
			if err != nil {
				return rpc.ErrDialRPC
			}

			fetchArgs := repo.FetchArgs{
				Remote: ctx.Config.Remote,
			}

			var fetchReply repo.FetchReply
			if err := client.Call("Repo.Fetch", &fetchArgs, &fetchReply); err != nil {
				return err
			}

			branch := ctx.Config.Branch
			heads := fetchReply.Repository.Heads()
			old := fetchReply.Repository.Branches[branch]
			new := ctx.Config.Branches[branch]

			pack, err := remote.BuildPack(c.Context, ctx.DAG, heads, old, new)
			if err != nil {
				return err
			}

			pushArgs := repo.PushArgs{
				Branch: branch,
				Pack:   pack,
				Remote: ctx.Config.Remote,
			}

			var pushReply repo.PushReply
			return client.Call("Repo.Push", &pushArgs, &pushReply)
		},
	}
}
