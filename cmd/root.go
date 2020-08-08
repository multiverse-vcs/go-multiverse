package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multi",
	Short: "Distributed version control.",
	Long:  `Distributed version control.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
