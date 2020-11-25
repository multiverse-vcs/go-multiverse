package storage

import (
	"path/filepath"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"github.com/spf13/afero"
)

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
