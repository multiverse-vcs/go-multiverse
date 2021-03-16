package remote

import (
	"context"
	"os"
	"path/filepath"

	badger "github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-path/resolver"

	"github.com/multiverse-vcs/go-multiverse/internal/key"
	"github.com/multiverse-vcs/go-multiverse/internal/p2p"
	"github.com/multiverse-vcs/go-multiverse/pkg/name"
)

// DotDir is the dot directory for the remote.
const DotDir = ".multiverse"

// Server implements the remote server.
type Server struct {
	// Config contains server settings.
	Config *Config
	// Keystore contains author keys.
	Keystore *key.Store
	// Namesys resolves named resources.
	Namesys *name.System
	// Peer manages peer services.
	Peer *p2p.Peer
	// Resolover is an ipfs path resolver.
	Resolver *resolver.Resolver
	// Root is the server root path.
	Root string
}

// NewServer returns a new remote server.
func NewServer(ctx context.Context, home string) (*Server, error) {
	root := filepath.Join(home, DotDir)
	if err := initServer(root); err != nil {
		return nil, err
	}

	config := NewConfig(root)
	if err := config.Read(); err != nil {
		return nil, err
	}

	priv, err := key.Decode(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	host, router, err := p2p.NewHost(ctx, priv, config.ListenAddresses)
	if err != nil {
		return nil, err
	}

	dpath := filepath.Join(root, "datastore")
	dopts := badger.DefaultOptions

	dstore, err := badger.NewDatastore(dpath, &dopts)
	if err != nil {
		return nil, err
	}

	peer, err := p2p.NewPeer(ctx, host, router, dstore)
	if err != nil {
		return nil, err
	}

	namesys, err := name.NewSystem(ctx, host, router, dstore)
	if err != nil {
		return nil, err
	}

	keystore, err := key.NewStore(root)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:   config,
		Keystore: keystore,
		Namesys:  namesys,
		Peer:     peer,
		Resolver: resolver.NewBasicResolver(peer.DAG),
		Root:     root,
	}, nil
}

// initServer initializes the remote server.
func initServer(root string) error {
	err := os.Mkdir(root, 0755)
	if os.IsExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	priv, err := key.Generate()
	if err != nil {
		return err
	}

	enc, err := key.Encode(priv)
	if err != nil {
		return err
	}

	config := NewConfig(root)
	config.PrivateKey = enc
	return config.Write()
}
