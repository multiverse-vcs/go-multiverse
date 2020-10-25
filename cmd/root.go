package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
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
