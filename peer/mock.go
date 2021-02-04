package peer

import (
	"context"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-ipfs-provider"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-path/resolver"
)

// Mock returns an offline peer with in memory storage.
func Mock(ctx context.Context, dstore datastore.Batching, config *Config) (*Client, error) {
	bstore := blockstore.NewBlockstore(dstore)
	exc := offline.Exchange(bstore)
	bserv := blockservice.New(bstore, exc)
	dag := merkledag.NewDAGService(bserv)
	resolv := resolver.NewBasicResolver(dag)
	system := provider.NewOfflineProvider()

	return &Client{
		DAGService: dag,
		config:     config,
		bstore:     bstore,
		dstore:     dstore,
		resolv:     resolv,
		system:     system,
	}, nil
}
