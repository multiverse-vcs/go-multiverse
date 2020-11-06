package cmd

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:          "branch",
	Short:        "Create, delete, or list branches.",
}

var branchCreateCmd = &cobra.Command{
	Use:          "add [name]",
	Short:        "Create a new branch.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeBranchCreate,
}

var branchDeleteCmd = &cobra.Command{
	Use:          "rm [name]",
	Short:        "Delete an existing branch.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeBranchDelete,
}

var branchListCmd = &cobra.Command{
	Use:          "ls",
	Short:        "Print all branches.",
	SilenceUsage: true,
	RunE:         executeBranchList,
}

func init() {
	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchCreateCmd)
	branchCmd.AddCommand(branchDeleteCmd)
	rootCmd.AddCommand(branchCmd)
}

func executeBranchCreate(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Open(cwd)
	if err != nil {
		return err
	}

	name := args[0]
	if _, ok := cfg.Branches[name]; ok {
		return fmt.Errorf("branch already exists")
	}

	cfg.Branches[name] = cfg.Base
	return cfg.Write()
}

func executeBranchDelete(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Open(cwd)
	if err != nil {
		return err
	}

	name := args[0]
	if _, ok := cfg.Branches[name]; ok {
		return fmt.Errorf("branch does not exists")
	}

	if name == cfg.Branch {
		return fmt.Errorf("cannot delete current branch")
	}

	delete(cfg.Branches, name)
	return cfg.Write()
}

func executeBranchList(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Open(cwd)
	if err != nil {
		return err
	}

	for name := range cfg.Branches {
		fmt.Printf("%s%s", colorYellow, name)

		if name == cfg.Branch {
			fmt.Printf(" (%sCURRENT%s)", colorGreen, colorYellow)
		}

		fmt.Printf("%s\n", colorReset)
	}

	return nil
}
