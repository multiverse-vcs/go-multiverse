package command

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/urfave/cli/v2"
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

			ctx, err := NewContext(cwd)
			if err != nil {
				return err
			}

			head, ok := ctx.Config.Repository.Branches[ctx.Config.Branch]
			if ok && head != ctx.Config.Index {
				return errors.New("index is behind head")
			}

			ignore, err := ctx.Ignore()
			if err != nil {
				return err
			}

			tree, err := fs.Add(c.Context, ctx.DAG, ctx.Root, ignore)
			if err != nil {
				return err
			}

			commit := object.NewCommit()
			commit.Tree = tree.Cid()
			commit.Message = c.String("message")

			commitID, err := object.AddCommit(c.Context, ctx.DAG, commit)
			if err != nil {
				return err
			}

			ctx.Config.Index = commitID
			ctx.Config.Repository.Branches[ctx.Config.Branch] = commitID
			return ctx.Config.Write()
		},
	}
}
