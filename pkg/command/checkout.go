package command

import (
	"errors"
	"os"

	cid "github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// NewCheckoutCommand returns a new cli command.
func NewCheckoutCommand() *cli.Command {
	return &cli.Command{
		Name:  "checkout",
		Usage: "Checkout committed files",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force checkout",
			},
			&cli.StringFlag{
				Name:    "commit",
				Aliases: []string{"c"},
				Usage:   "Checkout branch commit",
			},
			&cli.BoolFlag{
				Name:  "head",
				Usage: "Checkout branch head",
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

			branch := cc.Config.Branches[cc.Config.Branch]
			treeID := branch.Stash

			switch {
			case c.IsSet("head"):
				commit, err := object.GetCommit(c.Context, cc.DAG, branch.Head)
				if err != nil {
					return err
				}

				treeID = commit.Tree
			case c.IsSet("commit"):
				id, err := cid.Decode(c.String("commit"))
				if err != nil {
					return err
				}

				commit, err := object.GetCommit(c.Context, cc.DAG, id)
				if err != nil {
					return err
				}

				treeID = commit.Tree
			default:
				cli.ShowAppHelpAndExit(c, -1)
			}

			stash, err := fs.Add(c.Context, cc.DAG, cc.Root, context.DefaultIgnore)
			if err != nil {
				return err
			}

			status, err := dag.Status(c.Context, cc.DAG, stash, branch.Head)
			if err != nil {
				return err
			}

			if len(status) != 0 && !c.IsSet("force") {
				return errors.New("uncommitted changes")
			}

			tree, err := cc.DAG.Get(c.Context, treeID)
			if err != nil {
				return err
			}

			if err := fs.Write(c.Context, cc.DAG, cc.Root, tree); err != nil {
				return err
			}

			branch.Stash = treeID
			return cc.Config.Write()
		},
	}
}
