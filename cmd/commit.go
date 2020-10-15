package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var message string

var commitCmd = &cobra.Command{
	Use:          "commit",
	Short:        "Record changes to a repository.",
	SilenceUsage: true,
	RunE:         executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "description of changes")
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

	commit, err := c.Commit(cmd.Context(), message)
	if err != nil {
		return err
	}

	fmt.Println(commit.Cid().String())
	return nil
}
