package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/ipfs"
	"github.com/yondero/multiverse/repo"
)

var message string

var commitCmd = &cobra.Command{
	Use:          "commit",
	Short:        "Record changes in the local repository.",
	Long:         `Record changes in the local repository.`,
	SilenceUsage: true,
	RunE:         executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "description of changes")
}

func executeCommit(cmd *cobra.Command, args []string) error {
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

	c, err := r.Commit(ipfs, message)
	if err != nil {
		return err
	}

	fmt.Println(c.String())
	return nil
}
