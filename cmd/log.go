package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/core"
)

var logCmd = &cobra.Command{
	Use:          "log",
	Short:        "Print change history.",
	Long:         `Print change history.`,
	SilenceUsage: true,
	RunE:         executeLog,
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func executeLog(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	return core.Log(context.TODO(), cwd)
}
