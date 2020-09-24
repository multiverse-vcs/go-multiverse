package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/commit"
	"github.com/yondero/multiverse/repo"
)

var message string

var ignore = []string{repo.Config, ".git"}

var commitCmd = &cobra.Command{
	Use:          "commit",
	Short:        "Record changes to a repository.",
	Long:         `Record changes to a repository.`,
	SilenceUsage: true,
	RunE:         executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "description of changes")
}

func executeCommit(cmd *cobra.Command, args []string) error {
	api, err := httpapi.NewLocalApi()
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := repo.Open(cwd)
	if err != nil {
		return err
	}

	info, err := os.Stat(r.Path)
	if err != nil {
		return err
	}

	filter, err := files.NewFilter("", ignore, true)
	if err != nil {
		return err
	}

	tree, err := files.NewSerialFileWithFilter(r.Path, filter, info)
	if err != nil {
		return err
	}

	p, err := api.Unixfs().Add(context.TODO(), tree)
	if err != nil {
		return err
	}

	c := commit.Commit{Message: message, Tree: p.Root()}
	if r.Head.Defined() {
		c.Parents = append(c.Parents, r.Head)
	}

	node, err := c.Node()
	if err != nil {
		return err
	}

	if err := api.Dag().Pinning().Add(context.TODO(), node); err != nil {
		return err
	}

	r.Head = node.Cid()
	if err := r.Write(); err != nil {
		return err
	}

	fmt.Println(node.Cid().String())
	return nil
}
