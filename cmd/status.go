package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/config"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Print status of the local repo.",
	SilenceUsage: true,
	RunE:         executeStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func executeStatus(cmd *cobra.Command, args []string) error {
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

	diffs, err := c.Status(ctx, path.IpfsPath(config.Head), config.Path, config.Ignore...)
	if err != nil {
		return err
	}

	for _, d := range diffs {
		switch d.Type {
		case dagutils.Add:
			fmt.Printf("%sadded:   %s%s\n", ColorGreen, d.Path, ColorReset)
		case dagutils.Remove:
			fmt.Printf("%sremoved: %s%s\n", ColorRed, d.Path, ColorReset)
		case dagutils.Mod:
			fmt.Printf("%schanged: %s%s\n", ColorYellow, d.Path, ColorReset)
		}
	}

	return nil
}
