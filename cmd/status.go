package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/repo"
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

	r, err := repo.Open(cwd)
	if err != nil {
		return err
	}

	head, err := r.Branches.Head(r.Branch)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	tree, err := r.Tree()
	if err != nil {
		return err
	}

	changes, err := c.Status(ctx, path.IpfsPath(head), tree)
	if err != nil {
		return err
	}

	fmt.Printf("Tracking changes on branch %s:\n", r.Branch)
	fmt.Printf("  (all files are automatically considered for commit)\n")
	fmt.Printf("  (to stop tracking files add to '%s')\n", repo.IgnoreFile)

	for _, change := range changes {
		switch change.Type {
		case dagutils.Add:
			fmt.Printf("\t%snew file: %s%s\n", colorGreen, change.Path, colorReset)
		case dagutils.Remove:
			fmt.Printf("\t%sdeleted:  %s%s\n", colorRed, change.Path, colorReset)
		case dagutils.Mod:
			fmt.Printf("\t%smodified: %s%s\n", colorYellow, change.Path, colorReset)
		}
	}

	return nil
}
