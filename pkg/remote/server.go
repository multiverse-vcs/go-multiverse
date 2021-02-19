package remote

import (
	"context"
	"os"
	"path/filepath"

	badger "github.com/ipfs/go-ds-badger2"
	"github.com/multiverse-vcs/go-multiverse/pkg/name"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/multiverse-vcs/go-multiverse/pkg/p2p"
)

// DotDir is the dot directory for the remote.
const DotDir = ".multiversed"

// Server implements the remote server.
type Server struct {
	// Config contains server settings.
	Config *Config
	// Peer manages peer services.
	Peer *p2p.Peer
	// Namesys resolves named resources.
	Namesys *name.System
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

	key, err := p2p.DecodeKey(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	host, router, err := p2p.NewHost(ctx, key)
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

	authorID, err := object.AddAuthor(ctx, peer.DAG, config.Author)
	if err != nil {
		return nil, err
	}

	if err := namesys.Publish(ctx, key, authorID); err != nil {
		return nil, err
	}

	return &Server{
		Config:  config,
		Peer:    peer,
		Namesys: namesys,
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

	key, err := p2p.GenerateKey()
	if err != nil {
		return err
	}

	enc, err := p2p.EncodeKey(key)
	if err != nil {
		return err
	}

	config := NewConfig(root)
	config.PrivateKey = enc
	return config.Write()
}
