package cmd

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/repo"
	"github.com/spf13/cobra"
)

var branchDelete bool

var branchCmd = &cobra.Command{
	Use:          "branch",
	Short:        "List, create, or delete branches.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE:         executeBranch,
}

func init() {
	branchCmd.Flags().BoolVarP(&branchDelete, "delete", "d", false, "delete branch")
	rootCmd.AddCommand(branchCmd)
}

func executeBranch(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := repo.Open(cwd)
	if err != nil {
		return err
	}

	switch {
	case len(args) == 0:
		for name := range r.Branches {
			if name == r.Branch {
				fmt.Printf("* %s%s%s\n", colorGreen, name, colorReset)
			} else {
				fmt.Printf("%s%s%s\n", colorYellow, name, colorReset)
			}
		}

		return nil
	case branchDelete:
		if args[0] == r.Branch {
			return fmt.Errorf("cannot delete current branch")
		}

		if err := r.Branches.Remove(args[0]); err != nil {
			return nil
		}

		return r.Write()
	default:
		if err := r.Branches.Add(args[0], r.Base); err != nil {
			return nil
		}

		return r.Write()
	}
}
