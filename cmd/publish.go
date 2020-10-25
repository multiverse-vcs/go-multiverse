package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var (
	key string
	ref  string
)

var publishCmd = &cobra.Command{
	Use:          "publish [key] [ref]",
	Short:        "Announce a new version to peers.",
	Args:         cobra.RangeArgs(1, 2),
	SilenceUsage: true,
	RunE:         executePublish,
}

func init() {
	rootCmd.AddCommand(publishCmd)
}

func executePublish(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(cmd.Context(), config)
	if err != nil {
		return err
	}

	var p path.Path = path.IpfsPath(config.Head)
	if len(args) > 1 {
		p = path.New(args[1])
	}

	entry, err := c.Publish(cmd.Context(), args[0], p)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully published to %s\n", entry.Name())
	return nil
}
