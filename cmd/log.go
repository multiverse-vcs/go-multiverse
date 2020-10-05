package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ipfs/go-ipfs-http-client"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/core"
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
	addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		return err
	}

	api, err := httpapi.NewApi(addr)
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

	c := node.(*core.Commit)

	fmt.Printf("Commit: %s\n", node.Cid().String())
	fmt.Printf("Author: %s %s\n", c.Author.Name, c.Author.Email)
	fmt.Printf("Date:   %s\n", c.Author.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
	fmt.Printf("\n%s\n\n", c.Message)
	return nil
}
