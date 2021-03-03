package command

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// NewCommitCommand returns a new cli command.
func NewCommitCommand() *cli.Command {
	return &cli.Command{
		Name:  "commit",
		Usage: "Record a new version",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Usage:   "Description of the changes",
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

			tree, err := fs.Add(c.Context, cc.DAG, cc.Root, context.DefaultIgnore)
			if err != nil {
				return err
			}

			branch := cc.Config.Branches[cc.Config.Branch]
			diffs, err := dag.Status(c.Context, cc.DAG, tree, branch.Head)
			if err != nil {
				return err
			}

			if len(diffs) == 0 {
				return errors.New("no changes to commit")
			}

			commit := object.NewCommit()
			commit.Tree = tree.Cid()
			commit.Message = c.String("message")

			if branch.Head.Defined() {
				commit.Parents = append(commit.Parents, branch.Head)
			}

			commitID, err := object.AddCommit(c.Context, cc.DAG, commit)
			if err != nil {
				return err
			}

			branch.Head = commitID
			branch.Stash = tree.Cid()
			return cc.Config.Write()
		},
	}
}
