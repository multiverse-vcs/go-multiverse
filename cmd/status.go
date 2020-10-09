package cmd

import (
	"context"
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Print changes to working tree.",
	Long:         `Print changes to working tree.`,
	SilenceUsage: true,
	RunE:         executeStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func executeStatus(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(config)
	if err != nil {
		return err
	}

	return c.Status(context.TODO(), path.IpfsPath(config.Head))
}
