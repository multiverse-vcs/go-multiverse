package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yondero/go-ipld-multiverse"
	"github.com/yondero/go-multiverse/core"
)

var historyCmd = &cobra.Command{
	Use:          "history",
	Short:        "Print change history.",
	SilenceUsage: true,
	RunE:         executeHistory,
}

func init() {
	rootCmd.AddCommand(historyCmd)
}

func executeHistory(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := core.OpenConfig(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(cmd.Context(), config)
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

	return c.NewHistory(config.Head).ForEach(cmd.Context(), callback)
}
