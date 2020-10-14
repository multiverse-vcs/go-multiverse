package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Print status of the working tree.",
	Long:         `Print status of the working tree.`,
	SilenceUsage: true,
	RunE:         executeStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func executeStatus(cmd *cobra.Command, args []string) error {
	// TODO make background and cancel on interrupt
	ctx := context.TODO()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx, config)
	if err != nil {
		return err
	}

	return c.Status(ctx)
}
