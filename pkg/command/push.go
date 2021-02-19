package command

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
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

			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			client, err := remote.NewClient(home)
			if err != nil {
				return ErrDialRPC
			}

			fetchArgs := remote.FetchArgs{
				Remote: repo.Config.Remote,
			}

			var fetchReply remote.FetchReply
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

			pushArgs := remote.PushArgs{
				Remote: repo.Config.Remote,
				Branch: branch,
				Pack:   pack,
			}

			var pushReply remote.PushReply
			return client.Call("Remote.Push", &pushArgs, &pushReply)
		},
	}
}
