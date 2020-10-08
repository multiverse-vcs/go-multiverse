package cmd

import (
	"context"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var cloneCmd = &cobra.Command{
	Use:          "clone [remote] [local]",
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
	local, err := filepath.Abs(args[1])
	if err != nil {
		return err
	}

	config, err := core.InitConfig(local, cid.Cid{})
	if err != nil {
		return err
	}

	c, err := core.NewCore(config)
	if err != nil {
		return err
	}

	return c.Checkout(context.TODO(), path.New(args[0]))
}
