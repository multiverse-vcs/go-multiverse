package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
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

	diffs, err := c.Status(cmd.Context())
	if err != nil {
		return err
	}

	for _, d := range diffs {
		switch d.Type {
		case dagutils.Add:
			fmt.Printf("%s+ %s%s\n", ColorGreen, d.Path, ColorReset)
		case dagutils.Remove:
			fmt.Printf("%s- %s%s\n", ColorRed, d.Path, ColorReset)
		case dagutils.Mod:
			fmt.Printf("%s~ %s%s\n", ColorYellow, d.Path, ColorReset)
		}
	}

	return nil
}
