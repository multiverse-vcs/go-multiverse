package command

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
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

			repo, err := NewContext(cwd)
			if err != nil {
				return err
			}

			client, err := rpc.NewClient()
			if err != nil {
				return rpc.ErrDialRPC
			}

			fetchArgs := rpc.FetchArgs{
				Remote: repo.Config.Remote,
			}

			var fetchReply rpc.FetchReply
			if err := client.Call("Remote.Fetch", &fetchArgs, &fetchReply); err != nil {
				return err
			}

			branch := repo.Config.Branch
			heads := fetchReply.Repository.Heads()
			old := fetchReply.Repository.Branches[branch]
			new := repo.Config.Branches[branch]

			pack, err := remote.BuildPack(c.Context, repo.DAG, heads, old, new)
			if err != nil {
				return err
			}

			pushArgs := rpc.PushArgs{
				Branch: branch,
				Pack:   pack,
				Remote: repo.Config.Remote,
			}

			var pushReply rpc.PushReply
			return client.Call("Remote.Push", &pushArgs, &pushReply)
		},
	}
}
