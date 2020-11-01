package cmd

import (
	"os"

	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
	"github.com/yondero/go-multiverse/config"
)

var commitCmd = &cobra.Command{
	Use:          "commit [message]",
	Short:        "Record changes in the local repo.",
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

	config, err := config.Open(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	tree, err := c.WorkTree(ctx, config.Path)
	if err != nil {
		return err
	}

	opts := core.CommitOptions{
		Message:  args[0],
		Parents:  []cid.Cid{config.Head},
		Pin:      true,
		WorkTree: tree.Cid(),
	}

	commit, err := c.Commit(ctx, &opts)
	if err != nil {
		return err
	}

	config.Head = commit.Cid()
	return config.Write()
}
