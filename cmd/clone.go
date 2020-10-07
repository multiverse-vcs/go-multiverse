package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/core"
)

var cloneCmd = &cobra.Command{
	Use:          "clone [remote] [local]",
	Short:        "Copy an existing repository.",
	Long:         `Copy an existing repository.`,
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE:         executeClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func executeClone(cmd *cobra.Command, args []string) error {
	local, err := filepath.Abs(args[1])
	if err != nil {
		return err
	}

	config, err := core.Clone(context.TODO(), local, args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Repo cloned successfully to %s\n", config.Path)
	return nil
}
