package cmd

import (
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/repo"
	"github.com/multiverse-vcs/go-multiverse/util"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:          "init [ref]",
	Short:        "Create a new empty repo or copy an existing repo.",
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
	RunE:         executeInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func executeInit(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := repo.Init(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
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
	r.Branches[r.Branch] = commit.Cid()
	return r.Write()
}
