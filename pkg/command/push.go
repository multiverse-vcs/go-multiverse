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
		Name:      "push",
		Usage:     "Update a remote branch with local changes",
		ArgsUsage: "[remote] [branch]",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

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
			rname := c.Args().Get(0)
			bname := c.Args().Get(1)

			if !branch.Head.Defined() {
				return errors.New("nothing to push")
			}

			remote, ok := cc.Config.Remotes[rname]
			if !ok {
				return errors.New("remote does not exist")
			}

			args := repo.SearchArgs{
				Peer: remote.Peer,
				Name: remote.Name,
			}

			var reply repo.SearchReply
			if err := client.Call("Repo.Search", &args, &reply); err != nil {
				return err
			}

			refs := reply.Repository.Heads()
			head := reply.Repository.Branches[bname]

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
				Peer:   remote.Peer,
				Name:   remote.Name,
				Branch: bname,
				Data:   data.Bytes(),
			}

			return client.Call("Repo.Push", &pushArgs, nil)
		},
	}
}
