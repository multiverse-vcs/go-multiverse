package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/spf13/cobra"
	"github.com/yondero/go-multiverse/ipfs"
)

var daemonCmd = &cobra.Command{
	Use:          "daemon",
	Short:        "Runs a persistent Multiverse node.",
	Long:         `Runs a persistent Multiverse node.`,
	SilenceUsage: true,
	RunE:         executeDaemon,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

func executeDaemon(cmd *cobra.Command, args []string) error {
	plugins, err := ipfs.LoadPlugins()
	if err != nil {
		return err
	}

	node, err := ipfs.NewNode(context.TODO())
	if err != nil {
		return err
	}

	root, err := ipfs.RootPath()
	if err != nil {
		return err
	}

	cctx := commands.Context{
		ConfigRoot: root,
		Plugins:    plugins,
		ReqLog:     &commands.ReqLog{},
	}

	cctx.LoadConfig = func(path string) (*config.Config, error) {
		return node.Repo.Config()
	}

	cctx.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	go corehttp.ListenAndServe(node, ipfs.CommandsApiAddress, corehttp.CommandsOption(cctx))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	return nil
}
