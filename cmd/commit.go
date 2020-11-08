package cmd

import (
	"fmt"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/repo"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:          "commit [message]",
	Short:        "Record changes to the local repo.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func executeCommit(cmd *cobra.Command, args []string) error {
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

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	tree, err := r.Tree()
	if err != nil {
		return err
	}

	head, err := r.Branches.Head(r.Branch)
	if err != nil {
		return err
	}

	opts := core.CommitOptions{
		Message: args[0],
		Parents: []cid.Cid{head},
	}

	commit, err := c.Commit(ctx, tree, &opts)
	if err != nil {
		return err
	}

	defer fmt.Println(commit.Cid().String())

	r.Base = commit.Cid()
	r.Branches[r.Branch] = commit.Cid()
	return r.Write()
}
