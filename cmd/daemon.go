package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/spf13/cobra"
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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	plugins, err := loader.NewPluginLoader(filepath.Join(home, ".multi", "plugins"))
	if err != nil {
		return err
	}

	if err := plugins.Initialize(); err != nil {
		return err
	}

	if err := plugins.Inject(); err != nil {
		return err
	}

	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return err
	}

	path := filepath.Join(home, ".multi")
	if err := fsrepo.Init(path, cfg); err != nil {
		return err
	}

	repo, err := fsrepo.Open(path)
	if err != nil {
		return err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Repo: repo,
	}

	node, err := core.NewNode(context.TODO(), nodeOptions)
	if err != nil {
		return err
	}

	cctx := commands.Context{
		ConfigRoot: path,
		Plugins: plugins,
		ReqLog: &commands.ReqLog{},
		LoadConfig: func(path string) (*config.Config, error) {
			return node.Repo.Config()
		},
		ConstructNode: func() (*core.IpfsNode, error) {
			return node, nil
		},
	}

	opts := []corehttp.ServeOption{
		corehttp.GatewayOption(true, "/ipfs", "/ipns"),
		corehttp.WebUIOption,
		corehttp.CommandsOption(cctx),
	}

	return corehttp.ListenAndServe(node, "/ip4/127.0.0.1/tcp/5001", opts...)
}