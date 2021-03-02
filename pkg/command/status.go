package command

import (
	"fmt"
	"os"
	"sort"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/urfave/cli/v2"
)

// NewStatusCommand returns a new cli command.
func NewStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Print repository status",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			cc, err := context.New(cwd)
			if err != nil {
				return err
			}

			tree, err := cc.Tree(c.Context)
			if err != nil {
				return err
			}

			branch := cc.Config.Branches[cc.Config.Branch]

			diffs, err := dag.Status(c.Context, cc.DAG, tree, branch.Head)
			if err != nil {
				return err
			}

			paths := make([]string, 0)
			for path := range diffs {
				paths = append(paths, path)
			}
			sort.Strings(paths)

			fmt.Printf("Tracking changes on branch %s:\n", cc.Config.Branch)
			fmt.Printf("  (all files are automatically considered for commit)\n")
			fmt.Printf("  (to stop tracking files add rules to '%s')\n", context.IgnoreFile)

			for _, p := range paths {
				switch diffs[p] {
				case dagutils.Add:
					fmt.Printf("\tnew file: %s\n", p)
				case dagutils.Remove:
					fmt.Printf("\tdeleted:  %s\n", p)
				case dagutils.Mod:
					fmt.Printf("\tmodified: %s\n", p)
				}
			}

			return nil
		},
	}
}
