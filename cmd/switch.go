package cmd

import (
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:          "switch [name]",
	Short:        "Change to a different branch.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeSwitch,
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

func executeSwitch(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	name := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Open(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	id, ok := cfg.Branches[name]
	if !ok {
		return config.ErrBranchNotFound
	}

	commit, err := c.Reference(ctx, path.IpfsPath(id))
	if err != nil {
		return err
	}

	if err := c.Checkout(ctx, commit, cfg.Path); err != nil {
		return err
	}

	cfg.Base = id
	cfg.Branch = name
	return cfg.Write()
}
