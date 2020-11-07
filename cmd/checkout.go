package cmd

import (
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/repo"
	"github.com/multiverse-vcs/go-multiverse/util"
	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:          "checkout [ref]",
	Short:        "Copy changes from a commit to the local repo.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeCheckout,
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

func executeCheckout(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := repo.Open(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	commit, err := c.Reference(ctx, path.New(args[0]))
	if err != nil {
		return err
	}

	tree, err := c.Tree(ctx, commit)
	if err != nil {
		return err
	}

	if err := util.WriteTo(tree, r.Path); err != nil {
		return err
	}

	r.Base = commit.Cid()
	return r.Write()
}
