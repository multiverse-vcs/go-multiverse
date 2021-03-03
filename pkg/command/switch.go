package command

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
)

// NewSwitchCommand returns a new cli command.
func NewSwitchCommand() *cli.Command {
	return &cli.Command{
		Name:  "switch",
		Usage: "Change branches",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "keep",
				Aliases: []string{"k"},
				Usage:   "Keep working tree",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowAppHelpAndExit(c, -1)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			cc, err := context.New(cwd)
			if err != nil {
				return err
			}

			prev := cc.Config.Branch
			next := c.Args().Get(0)

			if prev == next {
				return errors.New("already on branch")
			}

			branch, ok := cc.Config.Branches[next]
			if !ok {
				return errors.New("branch does not exist")
			}

			stash, err := fs.Add(c.Context, cc.DAG, cc.Root, context.DefaultIgnore)
			if err != nil {
				return err
			}

			cc.Config.Branches[prev].Stash = stash.Cid()
			cc.Config.Branch = next

			if c.IsSet("keep") {
				return cc.Config.Write()
			}

			tree, err := cc.DAG.Get(c.Context, branch.Stash)
			if err != nil {
				return err
			}

			if err := fs.Write(c.Context, cc.DAG, cc.Root, tree); err != nil {
				return err
			}

			return cc.Config.Write()
		},
	}
}
