package cmd

import (
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/config"
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
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := config.Init(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
	}

	commit, err := c.Reference(ctx, path.New(args[0]))
	if err != nil {
		return err
	}

	if err := c.Checkout(ctx, commit, config.Path); err != nil {
		return err
	}

	config.Head = commit.Cid()
	return config.Write()
}
