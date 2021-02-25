package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	car "github.com/ipld/go-car"
	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
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

			ctx, err := context.New(cwd)
			if err != nil {
				return err
			}

			branch := c.String("branch")
			if branch == "" {
				branch = ctx.Config.Branch
			}

			head := ctx.Config.Repository.Branches[branch]
			if !head.Defined() {
				return errors.New("nothing to push")
			}

			fetchURL := fmt.Sprintf("http://%s/%s", remote.HttpAddr, ctx.Config.Remote)
			fetchRes, err := http.Get(fetchURL)
			if err != nil {
				return err
			}
			defer fetchRes.Body.Close()

			if fetchRes.StatusCode != http.StatusOK {
				return errors.New("fetch request failed")
			}

			var repo object.Repository
			if err := json.NewDecoder(fetchRes.Body).Decode(&repo); err != nil {
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

			pushURL := fmt.Sprintf("http://%s/%s/%s", remote.HttpAddr, ctx.Config.Remote, branch)
			pushRes, err := http.Post(pushURL, "application/octet-stream", &data)
			if err != nil {
				return err
			}
			defer pushRes.Body.Close()

			if pushRes.StatusCode != http.StatusCreated {
				return errors.New("push request failed")
			}

			return nil
		},
	}
}
