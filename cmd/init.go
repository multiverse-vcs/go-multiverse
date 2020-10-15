package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var initCmd = &cobra.Command{
	Use:          "init [ref]",
	Short:        "Create a new empty repo or copy an existing repo.",
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
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

	config, err := core.InitConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(cmd.Context(), config)
	if err != nil {
		return err
	}

	fmt.Printf("Repo initialized successfully at %s\n", config.Path)

	if len(args) > 0 {
		return c.Checkout(cmd.Context(), path.New(args[0]))
	}

	return nil
}
