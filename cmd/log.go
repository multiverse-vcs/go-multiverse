package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ipfs/go-ipfs-http-client"
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/commit"
	"github.com/yondero/multiverse/repo"
)

var logCmd = &cobra.Command{
	Use:          "log",
	Short:        "Print change history.",
	Long:         `Print change history.`,
	SilenceUsage: true,
	RunE:         executeLog,
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func executeLog(cmd *cobra.Command, args []string) error {
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

	if !r.Head.Defined() {
		return nil
	}

	node, err := api.Dag().Get(context.TODO(), r.Head)
	if err != nil {
		return err
	}

	c, err := commit.FromNode(node)
	if err != nil {
		return err
	}

	fmt.Printf("Commit: %s\n", node.Cid().String())
	fmt.Printf("Author: %s %s\n", c.Author.Name, c.Author.Email)
	fmt.Printf("Date:   %s\n", c.Author.When.Format("Mon Jan 2 15:04:05 2006 -0700"))
	fmt.Printf("\n%s\n\n", c.Message)
	return nil
}