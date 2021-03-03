package command

import (
	"fmt"
	"os"

	cid "github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/urfave/cli/v2"
)

// NewLogCommand returns a new cli command.
func NewLogCommand() *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "Print branch history",
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

			visit := func(id cid.Cid) bool {
				commit, err := object.GetCommit(c.Context, cc.DAG, id)
				if err != nil {
					return false
				}

				fmt.Printf("commit %s\n", id.String())
				fmt.Printf("Date:  %s\n", commit.Date.Format("Mon Jan 02 15:04:05 2006 -0700"))
				fmt.Printf("\n\t%s\n\n", commit.Message)
				return true
			}

			return dag.Walk(c.Context, cc.DAG, branch.Head, visit)
		},
	}
}
