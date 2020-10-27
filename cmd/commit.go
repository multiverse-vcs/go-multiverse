package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(cmd.Context(), config)
	if err != nil {
		return err
	}

	commit, err := c.Commit(cmd.Context(), args[0], config.Head)
	if err != nil {
		return err
	}

	config.Head = commit.Cid()
	return config.Write()
}
