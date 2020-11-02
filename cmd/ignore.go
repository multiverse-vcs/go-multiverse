package cmd

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/spf13/cobra"
)

var ignoreCmd = &cobra.Command{
	Use:          "ignore [pattern]",
	Short:        "Ignore changes to files matching pattern.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeIgnore,
}

func init() {
	rootCmd.AddCommand(ignoreCmd)
}

func executeIgnore(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := config.Open(cwd)
	if err != nil {
		return err
	}

	config.Ignore = append(config.Ignore, args[0])
	return config.Write()
}
