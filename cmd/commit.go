package cmd

import (
	"os"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:          "commit [message]",
	Short:        "Record changes to the local repo.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func executeCommit(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Open(cwd)
	if err != nil {
		return err
	}

	if err := cfg.Detached(); err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	tree, err := c.WorkTree(ctx, cfg.Path)
	if err != nil {
		return err
	}

	head, err := cfg.Head()
	if err != nil {
		return err
	}

	opts := core.CommitOptions{
		Message: args[0],
		Parents: []cid.Cid{head},
	}

	commit, err := c.Commit(ctx, tree.Cid(), &opts)
	if err != nil {
		return err
	}

	cfg.Base = commit.Cid()
	cfg.Branches[cfg.Branch] = commit.Cid()
	return cfg.Write()
}
