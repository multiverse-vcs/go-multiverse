package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var (
	name string
	ref  string
)

var publishCmd = &cobra.Command{
	Use:          "publish",
	Short:        "Announce a new version to peers.",
	Long:         `Announce a new version to peers.`,
	SilenceUsage: true,
	RunE:         executePublish,
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&name, "name", "n", "self", "name to publish under")
	publishCmd.Flags().StringVarP(&ref, "ref", "r", "", "reference to publish")
}

func executePublish(cmd *cobra.Command, args []string) error {
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

	var p path.Path = path.IpfsPath(config.Head)
	if len(ref) > 0 {
		p = path.New(ref)
	}

	entry, err := c.Publish(ctx, name, p)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully published to %s\n", entry.Name())
	return nil
}
