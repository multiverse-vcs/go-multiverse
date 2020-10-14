package cmd

import (
	"context"
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var checkoutCmd = &cobra.Command{
	Use:          "checkout [ref]",
	Short:        "Checkout a commit.",
	Long:         `Checkout a commit.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeCheckout,
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

func executeCheckout(cmd *cobra.Command, args []string) error {
	// TODO make background and cancel on interrupt
	ctx := context.TODO()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx, config)
	if err != nil {
		return err
	}

	return c.Checkout(ctx, path.New(args[0]))
}
