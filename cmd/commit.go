package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/repo"
)

var message string

var commitCmd = &cobra.Command{
	Use: "commit",
	Short: "Record changes in the local repository.",
	Long: `Record changes in the local repository.`,
	SilenceUsage: true,
	RunE: executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "description of changes")
}

func executeCommit(cmd *cobra.Command, args []string) error {
	r := repo.NewRepo()
	if err := r.Commit(message); err != nil {
		return err
	}

	return nil
}