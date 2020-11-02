package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/multiverse-vcs/go-ipld-multiverse"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/config"
)

var logCmd = &cobra.Command{
	Use:          "log",
	Short:        "Print change history.",
	SilenceUsage: true,
	RunE:         executeLog,
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func executeLog(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := config.Open(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	var callback core.HistoryCallback = func(commit *ipldmulti.Commit) error {
		fmt.Printf("%sCID:  %s%s\n", ColorCyan, commit.Cid().String(), ColorReset)
		fmt.Printf("Peer: %s\n", commit.PeerID.String())
		fmt.Printf("Date: %s\n", commit.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
		fmt.Printf("\n\t%s\n\n", commit.Message)
		return nil
	}

	return c.NewHistory(config.Head).ForEach(ctx, callback)
}
