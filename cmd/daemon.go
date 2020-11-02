package cmd

import (
	"os"
	"os/signal"

	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/multiformats/go-multiaddr/net"
	"github.com/multiverse-vcs/go-multiverse/ipfs"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:          "daemon",
	Short:        "Run a persistent Multiverse node.",
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

	node, err := ipfs.NewNode(cmd.Context())
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

	l, err := manet.Listen(ipfs.HttpApiAddress)
	if err != nil {
		return err
	}

	go corehttp.Serve(node, manet.NetListener(l), corehttp.CommandsOption(cctx))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	return nil
}
