package command

import (
	"fmt"
	"os"
	"sort"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
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

			ctx, err := context.New(cwd)
			if err != nil {
				return err
			}

			ignore, err := ctx.Ignore()
			if err != nil {
				return err
			}

			tree, err := fs.Add(c.Context, ctx.DAG, ctx.Root, ignore)
			if err != nil {
				return err
			}

			diffs, err := dag.Status(c.Context, ctx.DAG, ctx.Config.Index, tree)
			if err != nil {
				return err
			}

			paths := make([]string, 0)
			for path := range diffs {
				paths = append(paths, path)
			}
			sort.Strings(paths)

			fmt.Printf("Tracking changes on branch %s:\n", ctx.Config.Branch)
			fmt.Printf("  (all files are automatically considered for commit)\n")
			fmt.Printf("  (to stop tracking files add rules to '%s')\n", context.IgnoreFile)

			for _, p := range paths {
				switch diffs[p] {
				case dag.Add:
					fmt.Printf("\tnew file: %s\n", p)
				case dag.Remove:
					fmt.Printf("\tdeleted:  %s\n", p)
				case dag.Mod:
					fmt.Printf("\tmodified: %s\n", p)
				}
			}

			return nil
		},
	}
}
