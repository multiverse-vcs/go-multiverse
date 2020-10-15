package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var logCmd = &cobra.Command{
	Use:          "log",
	Short:        "Print change history.",
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

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(cmd.Context(), config)
	if err != nil {
		return err
	}

	return c.Log(cmd.Context(), config.Head)
}
