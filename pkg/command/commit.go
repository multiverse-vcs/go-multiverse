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

			repo, err := NewContext(cwd)
			if err != nil {
				return err
			}

			head, ok := repo.Config.Branches[repo.Config.Branch]
			if ok && head != repo.Config.Index {
				return errors.New("index is behind head")
			}

			ignore, err := repo.Ignore()
			if err != nil {
				return err
			}

			tree, err := fs.Add(c.Context, repo.DAG, repo.Root, ignore)
			if err != nil {
				return err
			}

			commit := object.NewCommit()
			commit.Tree = tree.Cid()
			commit.Message = c.String("message")

			commitID, err := object.AddCommit(c.Context, repo.DAG, commit)
			if err != nil {
				return err
			}

			repo.Config.Index = commitID
			repo.Config.Branches[repo.Config.Branch] = commitID
			return repo.Config.Write()
		},
	}
}
