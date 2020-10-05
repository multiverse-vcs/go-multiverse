// Package IPFS contains methods for running an IPFS node.
package ipfs

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

// RootPath returns the path to the root of the IPFS directory.
func RootPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".multi"), nil
}

// LoadPlugins initializes all plugins.
func LoadPlugins() (*loader.PluginLoader, error) {
	root, err := RootPath()
	if err != nil {
		return nil, err
	}

	plugins, err := loader.NewPluginLoader(filepath.Join(root, "plugins"))
	if err != nil {
		return nil, err
	}

	if err := plugins.Initialize(); err != nil {
		return nil, err
	}

	if err := plugins.Inject(); err != nil {
		return nil, err
	}

	return plugins, nil
}

// NewNode returns a new IPFS node.
func NewNode(ctx context.Context) (*core.IpfsNode, error) {
	root, err := RootPath()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return nil, err
	}

	if err := fsrepo.Init(root, cfg); err != nil {
		return nil, err
	}

	repo, err := fsrepo.Open(root)
	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online: true,
		Repo:   repo,
	}

	return core.NewNode(ctx, nodeOptions)
}