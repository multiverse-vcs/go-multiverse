package cmd

import (
	"context"
	"fmt"

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
	ipfs, err := ipfs.NewDefault(context.TODO())
	if err != nil {
		return err
	}

	c, err := repo.Commit(ipfs, message)
	if err != nil {
		return err
	}

	fmt.Println("Changes committed successfully")
	fmt.Println(c)
	return nil
}
