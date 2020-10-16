package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multi",
	Short: "Decentralized version control system.",
	Long: `Multiverse is a decentralized version control system
that enables peer-to-peer software development.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
