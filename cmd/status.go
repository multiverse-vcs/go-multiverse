package cmd

import (
	"fmt"

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

			for _, change := range changes {
				switch change.Type {
				case dagutils.Add:
					fmt.Printf("\t%snew file: %s%s\n", ColorGreen, change.Path, ColorReset)
				case dagutils.Remove:
					fmt.Printf("\t%sdeleted:  %s%s\n", ColorRed, change.Path, ColorReset)
				case dagutils.Mod:
					fmt.Printf("\t%smodified: %s%s\n", ColorYellow, change.Path, ColorReset)
				}
			}

			return nil
		},
	}
}
