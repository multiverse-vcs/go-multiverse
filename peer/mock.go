package peer

import (
	"context"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-ipfs-pinner/dspinner"
	"github.com/ipfs/go-ipfs-provider"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-path/resolver"
	"github.com/libp2p/go-libp2p-core/crypto"
)

// Mock returns an offline peer with in memory storage.
func Mock(ctx context.Context) (*Client, error) {
	dstore := datastore.NewMapDatastore()
	bstore := blockstore.NewBlockstore(dstore)
	exc := offline.Exchange(bstore)
	bserv := blockservice.New(bstore, exc)
	dag := merkledag.NewDAGService(bserv)
	resolv := resolver.NewBasicResolver(dag)
	system := provider.NewOfflineProvider()

	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	pinner, err := dspinner.New(ctx, dstore, dag)
	if err != nil {
		return nil, err
	}

	return &Client{
		DAGService: dag,
		Pinner:     pinner,
		priv:       priv,
		bstore:     bstore,
		dstore:     dstore,
		resolv:     resolv,
		system:     system,
	}, nil
}
