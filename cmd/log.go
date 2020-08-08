package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/repo"
)

var logCmd = &cobra.Command{
	Use:          "log",
	Short:        "Clone an existing Multiverse repository.",
	Long:         `Clone an existing Multiverse repository.`,
	SilenceUsage: true,
	RunE:         executeLog,
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func executeLog(cmd *cobra.Command, args []string) error {
	err := repo.Log()
	if err != nil {
		return err
	}

	return nil
}
