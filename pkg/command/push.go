package command

import (
	"bytes"
	"errors"
	"os"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	car "github.com/ipld/go-car"
	"github.com/multiverse-vcs/go-multiverse/pkg/http"
	"github.com/urfave/cli/v2"
)

// NewPushCommand returns a new cli command.
func NewPushCommand() *cli.Command {
	return &cli.Command{
		Name:  "push",
		Usage: "Update a remote repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "branch",
				Aliases: []string{"b"},
				Usage:   "Branch to push",
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

			branch := c.String("branch")
			if branch == "" {
				branch = ctx.Config.Branch
			}

			head := ctx.Config.Branches[branch]
			if !head.Defined() {
				return errors.New("nothing to push")
			}

			client := http.NewClient()
			repo, err := client.Fetch(ctx.Config.Remote)
			if err != nil {
				return err
			}

			// TODO use merge base to check if branch is valid

			refs := repo.Heads()
			walk := func(node ipld.Node) ([]*ipld.Link, error) {
				if refs.Has(node.Cid()) {
					return nil, nil
				}

				return node.Links(), nil
			}

			var data bytes.Buffer
			if err := car.WriteCarWithWalker(c.Context, ctx.DAG, []cid.Cid{head}, &data, walk); err != nil {
				return err
			}

			return client.Push(ctx.Config.Remote, branch, data.Bytes())
		},
	}
}
