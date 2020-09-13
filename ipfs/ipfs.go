package ipfs

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-ds-badger"
	"github.com/ipfs/go-ipfs-blockstore"
)

// Node is a wrapper around IPFS services.
type Node struct {
	Blocks blockservice.BlockService
}

// NewNode creates an IPFS node.
func NewNode(ctx context.Context) (*Node, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".multiverse", "datastore")
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, err
	}

	store, err := badger.NewDatastore(path, &badger.DefaultOptions)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(store)
	bstore = blockstore.NewIdStore(bstore)

	opts := blockstore.DefaultCacheOpts()
	bstore, err = blockstore.CachedBlockstore(ctx, bstore, opts)
	if err != nil {
		return nil, err
	}

	bservice := blockservice.New(bstore, nil)
	return &Node{Blocks: bservice}, nil
}