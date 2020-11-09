package cmd

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/port"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:          "import [type]",
	Short:        "Import a repo from an external VCS.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func executeImport(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	var importer port.Importer
	switch args[0] {
	case "git":
		importer = port.NewGitImporter(c)
	default:
		return fmt.Errorf("invalid import type")
	}

	return importer.Import(ctx, cwd)
}
