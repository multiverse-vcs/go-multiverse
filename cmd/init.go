package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/repo"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Create a new empty repo.",
	Long:         `Create a new empty repo.`,
	SilenceUsage: true,
	RunE:         executeInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func executeInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = repo.Init(cwd)
	if err != nil {
		return err
	}

	fmt.Println("Repo initialized successfully!")
	return nil
}
