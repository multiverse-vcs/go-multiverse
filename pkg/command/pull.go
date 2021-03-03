package command

import (
	"bytes"
	"errors"
	"os"

	cid "github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/multiverse-vcs/go-multiverse/pkg/merge"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
)

// NewPullCommand returns a new cli command.
func NewPullCommand() *cli.Command {
	return &cli.Command{
		Name:  "pull",
		Usage: "Update a local branch with remote changes",
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
			source := cc.Config.Branch

			stash, err := fs.Add(c.Context, cc.DAG, cc.Root, context.DefaultIgnore)
			if err != nil {
				return err
			}

			status, err := dag.Status(c.Context, cc.DAG, stash, branch.Head)
			if err != nil {
				return err
			}

			if len(status) != 0 {
				return errors.New("uncommitted changes")
			}

			if c.IsSet("remote") {
				remote = c.String("remote")
			}

			if c.IsSet("branch") {
				source = c.String("branch")
			}

			if alias, ok := cc.Config.Remotes[remote]; ok {
				remote = alias
			}

			var refs []cid.Cid
			for _, b := range cc.Config.Branches {
				if !b.Head.Defined() {
					continue
				}

				refs = append(refs, b.Head)
			}

			args := repo.PullArgs{
				Remote: remote,
				Branch: source,
				Refs:   refs,
			}

			var reply repo.PullReply
			if err := client.Call("Repo.Pull", &args, &reply); err != nil {
				return err
			}

			root, err := dag.ReadCar(cc.Blocks, bytes.NewReader(reply.Data))
			if err != nil {
				return err
			}

			base, err := merge.Base(c.Context, cc.DAG, branch.Head, root)
			if err != nil {
				return err
			}

			tree, err := merge.Tree(c.Context, cc.DAG, base, branch.Head, root)
			if err != nil {
				return err
			}

			if err := fs.Write(c.Context, cc.DAG, cc.Root, tree); err != nil {
				return err
			}

			branch.Head = root
			branch.Stash = tree.Cid()
			return cc.Config.Write()
		},
	}
}
