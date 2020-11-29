package storage

import (
	"os"
	"path/filepath"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/spf13/afero"
)

// InitOsStore initializes a store that is backed by the operating system.
func InitOsStore(root string) (*Store, error) {
	path := filepath.Join(root, DotDir)
	if err := os.Mkdir(path, 0755); err != nil {
		return nil, err
	}

	store, err := NewOsStore(root)
	if err != nil {
		return nil, err
	}

	priv, _, err := crypto.GenerateKeyPair(KeyType, -1)
	if err != nil {
		return nil, err
	}

	if err := store.WriteConfig(config.Default()); err != nil {
		return nil, err
	}

	if err := store.WriteKey(priv); err != nil {
		return nil, err
	}

	return store, nil
}

// NewOsStore returns a store that is backed by the operating system.
func NewOsStore(root string) (*Store, error) {
	cwd := afero.NewBasePathFs(afero.NewOsFs(), root)
	dot := afero.NewBasePathFs(cwd, DotDir)

	opts := badger.DefaultOptions
	path := filepath.Join(root, DotDir, DataDir)

	dstore, err := badger.NewDatastore(path, &opts)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(dstore)
	exc := offline.Exchange(bstore)

	bserv := blockservice.New(bstore, exc)
	dag := merkledag.NewDAGService(bserv)

	return &Store{
		Dag:    dag,
		Dot:    dot,
		Cwd:    cwd,
		bstore: bstore,
	}, nil
}
