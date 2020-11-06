package cmd

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/spf13/cobra"
)

var branchDelete bool

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List, create, or delete branches.",
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

	cfg, err := config.Open(cwd)
	if err != nil {
		return err
	}

	switch {
	case len(args) == 0:
		return executeBranchList(cfg)
	case branchDelete:
		return executeBranchDelete(cfg, args[0])
	default:
		return executeBranchCreate(cfg, args[0])
	}

	return nil
}

func executeBranchCreate(cfg *config.Config, name string) error {
	if err := cfg.Branches.Add(name, cfg.Base); err != nil {
		return nil
	}

	return cfg.Write()
}

func executeBranchDelete(cfg *config.Config, name string) error {
	if name == cfg.Branch {
		return fmt.Errorf("cannot delete current branch")
	}

	if err := cfg.Branches.Remove(name); err != nil {
		return nil
	}

	return cfg.Write()
}

func executeBranchList(cfg *config.Config) error {
	for name := range cfg.Branches {
		if name == cfg.Branch {
			fmt.Printf("* %s%s%s\n", colorGreen, name, colorReset)
		} else {
			fmt.Printf("%s%s%s\n", colorYellow, name, colorReset)
		}
	}

	return nil
}