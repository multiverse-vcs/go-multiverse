package cmd

import (
	"github.com/multiverse-vcs/go-multiverse/util"
	"github.com/spf13/cobra"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
)

var rootCmd = &cobra.Command{
	Use:   "multi",
	Short: "Decentralized version control system.",
	Long: `Multiverse is a decentralized version control system
that enables peer-to-peer software development.`,
}

// Execute runs the root command.
func Execute() error {
	if err := util.SetFileLimit(2048, 8192); err != nil {
		return err
	}

	return rootCmd.Execute()
}
