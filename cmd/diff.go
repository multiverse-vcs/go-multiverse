package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var diffCmd = &cobra.Command{
	Use:          "diff",
	Short:        "Print changes to the working tree.",
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

	c, err := core.NewCore(cmd.Context(), config)
	if err != nil {
		return err
	}

	return c.Diff(cmd.Context())
}
