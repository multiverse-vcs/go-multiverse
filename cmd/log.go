package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/ipfs"
	"github.com/yondero/multiverse/repo"
)

var logCmd = &cobra.Command{
	Use:          "log",
	Short:        "Print change history.",
	Long:         `Print change history.`,
	SilenceUsage: true,
	RunE:         executeLog,
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func executeLog(cmd *cobra.Command, args []string) error {
	ipfs, err := ipfs.NewNode(context.TODO())
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := repo.Open(cwd)
	if err != nil {
		return err
	}

	id, err := r.Head()
	if err != nil {
		return err
	}

	return r.Log(ipfs, id)
}
