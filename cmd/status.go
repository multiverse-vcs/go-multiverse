package cmd

import (
	"fmt"
	"sort"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/urfave/cli/v2"
)

// NewStatusCommand returns a new status command.
func NewStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Print repo status",
		Action: func(c *cli.Context) error {
			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			cfg, err := store.ReadConfig()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			changes, err := core.Status(c.Context, store, cfg.Head())
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("Tracking changes on branch %s:\n", cfg.Branch)
			fmt.Printf("  (all files are automatically considered for commit)\n")
			fmt.Printf("  (to stop tracking files add rules to '%s')\n", core.IgnoreFile)

			set := make(map[string]dagutils.ChangeType)
			for _, change := range changes {
				if _, ok := set[change.Path]; ok {
					set[change.Path] = dagutils.Mod
				} else {
					set[change.Path] = change.Type
				}
			}

			paths := make([]string, 0)
			for path := range set {
				paths = append(paths, path)
			}
			sort.Strings(paths)

			for _, p := range paths {
				switch set[p] {
				case dagutils.Add:
					fmt.Printf("\t%snew file: %s%s\n", ColorGreen, p, ColorReset)
				case dagutils.Remove:
					fmt.Printf("\t%sdeleted:  %s%s\n", ColorRed, p, ColorReset)
				case dagutils.Mod:
					fmt.Printf("\t%smodified: %s%s\n", ColorRed, p, ColorReset)
				}
			}

			return nil
		},
	}
}
