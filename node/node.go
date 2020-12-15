package node

import (
	"path/filepath"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
)

// DataDir is the name of the data directory.
const DataDir = "datastore"

// Node contains dag and libp2p services.
type Node struct {
	// Dag is the merkledag interface.
	Dag ipld.DAGService
	// Host is the optional p2p host.
	Host host.Host

	router routing.Routing
	bstore blockstore.Blockstore
}

// NewNode returns a new node.
func NewNode(root string) (*Node, error) {
	path := filepath.Join(root, DataDir)
	opts := badger.DefaultOptions

	dstore, err := badger.NewDatastore(path, &opts)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(dstore)
	exc := offline.Exchange(bstore)
	bserv := blockservice.New(bstore, exc)

	return &Node{
		Dag:    merkledag.NewDAGService(bserv),
		bstore: bstore,
	}, nil
}
