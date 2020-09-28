package cmd

import (
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/importer"
)

var importCmd = &cobra.Command{
	Use:          "import [local]",
	Short:        "Import an existing git repo.",
	Long:         `Import an existing git repo.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func executeImport(cmd *cobra.Command, args []string) error {
	api, err := httpapi.NewLocalApi()
	if err != nil {
		return err
	}

	return importer.NewGitImporter(api).Import(args[0])
}