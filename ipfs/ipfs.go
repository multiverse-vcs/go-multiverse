package ipfs

import (
	"context"
	"io/ioutil"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/go-ipfs/plugin/loader"
)

func initPlugins(path string) error {
	plugins, err := loader.NewPluginLoader(path)
	if err != nil {
		return err
	}

	if err := plugins.Initialize(); err != nil {
		return err
	}

	return plugins.Inject()
}

func initRepo(repoPath string) (repo.Repo, error) {
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return nil, err
	}

	if err = fsrepo.Init(repoPath, cfg); err != nil {
		return nil, err
	}

	return fsrepo.Open(repoPath)
}

// NewDefault creates a default IPFS node.
func NewDefault(ctx context.Context) (*core.IpfsNode, error) {
	repoPath, err := fsrepo.BestKnownPath()
	if err != nil {
		return nil, err
	}

	pluginsPath := filepath.Join(repoPath, "plugins")
	if err = initPlugins(pluginsPath); err != nil {
		return nil, err
	}

	repo, err := initRepo(repoPath)
	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online: true,
		Routing: libp2p.DHTOption,
		Repo: repo,
	}

	return core.NewNode(ctx, nodeOptions)
}

// NewEphemeral creates an ephemeral IPFS node.
func NewEphemeral(ctx context.Context) (*core.IpfsNode, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return nil, err
	}

	pluginsPath := filepath.Join(repoPath, "plugins")
	if err = initPlugins(pluginsPath); err != nil {
		return nil, err
	}

	repo, err := initRepo(repoPath)
	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online: true,
		Routing: libp2p.DHTOption,
		Repo: repo,
	}

	return core.NewNode(ctx, nodeOptions)
}