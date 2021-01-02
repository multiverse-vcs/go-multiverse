package node

import (
	"context"
	"time"

	"github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-provider"
	"github.com/ipfs/go-ipfs-provider/queue"
	"github.com/ipfs/go-ipfs-provider/simple"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-path/resolver"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

const (
	// ReprovideInterval is the time between reprovides.
	ReprovideInterval = 12 * time.Hour
	// QueueName is the name for the provider queue.
	QueueName = "repro"
)

// Init initializes and returns a new node.
func Init(ctx context.Context, dstore datastore.Batching) (*Node, error) {
	key, err := p2p.GenerateKey()
	if err != nil {
		return nil, err
	}

	host, router, err := p2p.NewHost(ctx, key)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(dstore)
	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, bstore)
	bserv := blockservice.New(bstore, exc)

	for _, info := range dht.GetDefaultBootstrapPeerAddrInfos() {
		go host.Connect(ctx, info)
	}

	if err := p2p.Discovery(ctx, host); err != nil {
		return nil, err
	}

	queue, err := queue.NewQueue(ctx, QueueName, dstore)
	if err != nil {
		return nil, err
	}

	prov := simple.NewProvider(ctx, queue, router)
	keys := simple.NewBlockstoreProvider(bstore)
	reprov := simple.NewReprovider(ctx, ReprovideInterval, router, keys)

	system := provider.NewSystem(prov, reprov)
	system.Run()

	dag := merkledag.NewDAGService(bserv)
	resolv := resolver.NewBasicResolver(dag)

	return &Node{
		DAGService: dag,
		host:       host,
		bstore:     bstore,
		dstore:     dstore,
		resolv:     resolv,
		router:     router,
		system:     system,
	}, nil
}
