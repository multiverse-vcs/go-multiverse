package cmd

import (
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/core"
)

var mergeCmd = &cobra.Command{
	Use:          "merge [ref]",
	Short:        "Merge changes from a peer into the local repo.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeMerge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}

func executeMerge(cmd *cobra.Command, args []string) error {
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

	commit, err := c.Merge(cmd.Context(), path.New(args[0]))
	if err != nil {
		return err
	}

	config.Head = commit.Cid()
	return config.Write()
}
