// Package IPFS contains methods for running an IPFS node.
package ipfs

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/multiformats/go-multiaddr"
	"github.com/yondero/go-ipld-multiverse"
)

// HttpApiAddress is the multiaddress of the commands API.
var HttpApiAddress = multiaddr.StringCast("/ip4/127.0.0.1/tcp/5001")

func init() {
	cid.Codecs["multi-commit"] = ipldmulti.CommitCodec
	cid.CodecToStr[ipldmulti.CommitCodec] = "multi-commit"
	format.Register(ipldmulti.CommitCodec, ipldmulti.DecodeCommit)
}

// RootPath returns the path to the root of the IPFS directory.
func RootPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".multi"), nil
}

// LoadPlugins loads and initializes plugins.
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

// Initializes a node with default configuration settings.
func Initialize(root string) error {
	if fsrepo.IsInitialized(root) {
		return nil
	}

	cfg, err := config.Init(os.Stdout, 2048)
	if err != nil {
		return err
	}

	return fsrepo.Init(root, cfg)
}

// NewNode returns a new IPFS node.
func NewNode(ctx context.Context) (*core.IpfsNode, error) {
	root, err := RootPath()
	if err != nil {
		return nil, err
	}

	if err := Initialize(root); err != nil {
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

// NewApi returns a new IPFS core API.
// Falls back to an HTTP API if repo is locked.
func NewApi(ctx context.Context) (iface.CoreAPI, error) {
	root, err := RootPath()
	if err != nil {
		return nil, err
	}

	locked, err := fsrepo.LockedByOtherProcess(root)
	if err != nil {
		return nil, err
	}

	if locked {
		return httpapi.NewApi(HttpApiAddress)
	}

	_, err = LoadPlugins()
	if err != nil {
		return nil, err
	}

	node, err := NewNode(ctx)
	if err != nil {
		return nil, err
	}

	return coreapi.NewCoreAPI(node)
}
