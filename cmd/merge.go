package cmd

import (
	"os"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/repo"
	"github.com/multiverse-vcs/go-multiverse/util"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:          "merge [ref] [message]",
	Short:        "Merge changes from a peer into the local repo.",
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE:         executeMerge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}

func executeMerge(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := repo.Open(cwd)
	if err != nil {
		return err
	}

	if err := r.Detached(); err != nil {
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

	local, err := c.Reference(ctx, path.IpfsPath(head))
	if err != nil {
		return err
	}

	remote, err := c.Reference(ctx, path.New(args[0]))
	if err != nil {
		return err
	}

	merge, err := c.Merge(ctx, local, remote)
	if err != nil {
		return err
	}

	opts := core.CommitOptions{
		Message: args[1],
		Parents: []cid.Cid{local.Cid(), remote.Cid()},
	}

	commit, err := c.Commit(ctx, merge, &opts)
	if err != nil {
		return err
	}

	if err := util.WriteTo(merge, r.Path); err != nil {
		return err
	}

	r.Base = commit.Cid()
	r.Branches[r.Branch] = commit.Cid()
	return r.Write()
}
