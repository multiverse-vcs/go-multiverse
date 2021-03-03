package p2p

import (
	"context"
	"time"

	bitswap "github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	blockservice "github.com/ipfs/go-blockservice"
	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	provider "github.com/ipfs/go-ipfs-provider"
	"github.com/ipfs/go-ipfs-provider/queue"
	"github.com/ipfs/go-ipfs-provider/simple"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
)

const (
	// ReprovideInterval is the time between reprovides.
	ReprovideInterval = 12 * time.Hour
	// QueueName is the name for the provider queue.
	QueueName = "repro"
)

// Peer implements p2p services.
type Peer struct {
	// Blocks is the ipfs blockstore.
	Blocks blockstore.Blockstore
	// DAG implements ipld DAGService.
	DAG ipld.DAGService
	// Host is the libp2p host.
	Host host.Host
	// Router is the libp2p router.
	Router routing.Routing
}

// New returns a new peer using the given host, router, and datstore.
func NewPeer(ctx context.Context, host host.Host, router routing.Routing, dstore datastore.Batching) (*Peer, error) {
	bstore := blockstore.NewBlockstore(dstore)
	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, bstore)
	bserv := blockservice.New(bstore, exc)

	queue, err := queue.NewQueue(ctx, QueueName, dstore)
	if err != nil {
		return nil, err
	}

	prov := simple.NewProvider(ctx, queue, router)
	bspr := simple.NewBlockstoreProvider(bstore)
	repr := simple.NewReprovider(ctx, ReprovideInterval, router, bspr)

	sys := provider.NewSystem(prov, repr)
	sys.Run()

	if err := Discovery(ctx, host); err != nil {
		return nil, err
	}

	if err := Bootstrap(ctx, host, router); err != nil {
		return nil, err
	}

	return &Peer{
		Blocks: bstore,
		DAG:    merkledag.NewDAGService(bserv),
		Host:   host,
		Router: router,
	}, nil
}
