package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multi",
	Short: "Distributed version control.",
	Long: `Distributed version control.`,
}

func Execute() error {
	return rootCmd.Execute();
}