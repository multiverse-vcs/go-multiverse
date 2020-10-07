package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var message string

var commitCmd = &cobra.Command{
	Use:          "commit",
	Short:        "Record changes to a repository.",
	Long:         `Record changes to a repository.`,
	SilenceUsage: true,
	RunE:         executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "description of changes")
}

func executeCommit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c, err := core.Commit(context.TODO(), cwd, message)
	if err != nil {
		return err
	}

	fmt.Println(c.Cid().String())
	return nil
}
