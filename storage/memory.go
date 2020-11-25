package storage

import (
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"github.com/spf13/afero"
)

// NewMemoryStore returns an in memory only store.
func NewMemoryStore() (*Store, error) {
	cwd := afero.NewMemMapFs()
	dot := afero.NewBasePathFs(cwd, DotDir)

	opts := badger.DefaultOptions
	opts.Options = opts.WithInMemory(true)

	dstore, err := badger.NewDatastore("", &opts)
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
