package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/spf13/cobra"
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

	cfg, err := config.Open(cwd)
	if err != nil {
		return err
	}

	head, err := cfg.Head()
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	changes, err := c.Status(ctx, path.IpfsPath(head), cfg.Path)
	if err != nil {
		return err
	}

	for _, change := range changes {
		switch change.Type {
		case dagutils.Add:
			fmt.Printf("%sAdd:    %s%s\n", colorGreen, change.Path, colorReset)
		case dagutils.Remove:
			fmt.Printf("%sRemove: %s%s\n", colorRed, change.Path, colorReset)
		case dagutils.Mod:
			fmt.Printf("%sModify: %s%s\n", colorYellow, change.Path, colorReset)
		}
	}

	return nil
}
