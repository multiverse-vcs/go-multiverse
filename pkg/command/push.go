package command

import (
	"bytes"
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/merge"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
)

// NewPushCommand returns a new cli command.
func NewPushCommand() *cli.Command {
	return &cli.Command{
		Name:  "push",
		Usage: "Update a remote branch with local changes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "remote",
				Aliases: []string{"r"},
				Usage:   "Remote repository path",
			},
			&cli.StringFlag{
				Name:    "branch",
				Aliases: []string{"b"},
				Usage:   "Remote branch name",
			},
		},
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			cc, err := context.New(cwd)
			if err != nil {
				return err
			}

			client, err := rpc.NewClient()
			if err != nil {
				return cli.Exit(rpc.DialErrMsg, -1)
			}

			branch := cc.Config.Branches[cc.Config.Branch]
			remote := branch.Remote
			target := cc.Config.Branch

			if !branch.Head.Defined() {
				return errors.New("nothing to push")
			}

			if c.IsSet("remote") {
				remote = c.String("remote")
			}

			if c.IsSet("branch") {
				target = c.String("branch")
			}

			if alias, ok := cc.Config.Remotes[remote]; ok {
				remote = alias
			}

			args := repo.SearchArgs{
				Remote: remote,
			}

			var reply repo.SearchReply
			if err := client.Call("Repo.Search", &args, &reply); err != nil {
				return err
			}

			refs := reply.Repository.Heads()
			head := reply.Repository.Branches[target]

			base, err := merge.Base(c.Context, cc.DAG, head, branch.Head)
			if err != nil {
				return err
			}

			if base != head {
				return errors.New("branches are non-divergent")
			}

			var data bytes.Buffer
			if err := dag.WriteCar(c.Context, cc.DAG, branch.Head, refs, &data); err != nil {
				return err
			}

			pushArgs := repo.PushArgs{
				Remote: remote,
				Branch: target,
				Data:   data.Bytes(),
			}

			return client.Call("Repo.Push", &pushArgs, nil)
		},
	}
}
