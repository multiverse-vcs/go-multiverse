package cmd

import (
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
	"github.com/yondero/go-multiverse/config"
)

var checkoutCmd = &cobra.Command{
	Use:          "checkout [ref]",
	Short:        "Checkout files from a different commit.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeCheckout,
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

func executeCheckout(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := config.Open(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	commit, err := c.Checkout(ctx, path.New(args[0]), config.Path)
	if err != nil {
		return err
	}

	config.Head = commit.Cid()
	return config.Write()
}
