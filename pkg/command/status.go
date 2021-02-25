package command

import (
	"fmt"
	"os"
	"sort"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
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

			ctx, err := NewContext(cwd)
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

			commit, err := object.GetCommit(c.Context, ctx.DAG, ctx.Config.Index)
			if err != nil {
				return err
			}

			index, err := ctx.DAG.Get(c.Context, commit.Tree)
			if err != nil {
				return err
			}

			changes, err := dagutils.Diff(c.Context, ctx.DAG, index, tree)
			if err != nil {
				return err
			}

			diffs := make(map[string]dagutils.ChangeType)
			for _, change := range changes {
				if _, ok := diffs[change.Path]; ok {
					diffs[change.Path] = dagutils.Mod
				} else if change.Path != "" {
					diffs[change.Path] = change.Type
				}
			}

			paths := make([]string, 0)
			for path := range diffs {
				paths = append(paths, path)
			}
			sort.Strings(paths)

			fmt.Printf("Tracking changes on branch %s:\n", ctx.Config.Branch)
			fmt.Printf("  (all files are automatically considered for commit)\n")
			fmt.Printf("  (to stop tracking files add rules to '%s')\n", IgnoreFile)

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
