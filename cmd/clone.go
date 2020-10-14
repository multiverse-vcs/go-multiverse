package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var cloneCmd = &cobra.Command{
	Use:          "clone [ref] [path]",
	Short:        "Copy an existing repository.",
	Long:         `Copy an existing repository.`,
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE:         executeClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func executeClone(cmd *cobra.Command, args []string) error {
	// TODO make background and cancel on interrupt
	ctx := context.TODO()

	local, err := filepath.Abs(args[1])
	if err != nil {
		return err
	}

	if err := os.Mkdir(local, 0777); err != nil {
		return err
	}

	config, err := core.InitConfig(local, cid.Cid{})
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx, config)
	if err != nil {
		return err
	}

	return c.Checkout(ctx, path.New(args[0]))
}
