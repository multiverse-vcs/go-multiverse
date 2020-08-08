package cmd

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/repo"
)

var cloneCmd = &cobra.Command{
	Use:          "clone [cid] [path]",
	Short:        "Clone an existing Multiverse repository.",
	Long:         `Clone an existing Multiverse repository.`,
	Args:         cobra.MinimumNArgs(2),
	SilenceUsage: true,
	RunE:         executeClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func executeClone(cmd *cobra.Command, args []string) error {
	id, err := cid.Parse(args[0])
	if err != nil {
		return err
	}

	dir, err := repo.Clone(id, args[1])
	if err != nil {
		return err
	}

	fmt.Println("Repo cloned at", dir)
	return nil
}
