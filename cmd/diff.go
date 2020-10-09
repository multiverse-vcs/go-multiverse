package cmd

import (
	"context"
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var diffCmd = &cobra.Command{
	Use:          "diff [remoteA] [remoteB]",
	Short:        "Print changes between commits.",
	Long:         `Print changes between commits.`,
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE:         executeDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func executeDiff(cmd *cobra.Command, args []string) error {
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

	return c.Diff(context.TODO(), path.New(args[0]), path.New(args[1]))
}
