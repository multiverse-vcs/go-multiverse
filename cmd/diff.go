package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var diffCmd = &cobra.Command{
	Use:          "diff",
	Short:        "Prints changes to files in the working tree.",
	Long:         `Prints changes to files in the working tree.`,
	SilenceUsage: true,
	RunE:         executeDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func executeDiff(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(config)
	if err != nil {
		return err
	}

	return c.Diff(context.TODO())
}
