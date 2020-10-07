package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/core"
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

	config, err = core.InitConfig(cwd, cid.Cid{})
	if err != nil {
		return err
	}

	fmt.Printf("Repo initialized successfully at %s\n", config.Path)
	return nil
}
