package cmd

import (
	"os"
	
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multi",
	Short: "Distributed version control.",
	Long: `Distributed version control.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}