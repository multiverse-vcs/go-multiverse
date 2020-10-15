package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Print status of the local repo.",
	SilenceUsage: true,
	RunE:         executeStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func executeStatus(cmd *cobra.Command, args []string) error {
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

	return c.Status(cmd.Context())
}
