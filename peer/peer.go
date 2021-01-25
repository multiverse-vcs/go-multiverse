// Package peer implements a peer client.
package peer

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
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

const (
	// ReprovideInterval is the time between reprovides.
	ReprovideInterval = 12 * time.Hour
	// QueueName is the name for the provider queue.
	QueueName = "repro"
)

// Client manages peer services.
type Client struct {
	ipld.DAGService
	host   host.Host
	priv   crypto.PrivKey
	bstore blockstore.Blockstore
	dstore datastore.Batching
	resolv *resolver.Resolver
	router routing.Routing
	system provider.System
}

// New returns a peer with a p2p host and persistent storage.
func New(ctx context.Context, dstore datastore.Batching, priv crypto.PrivKey) (*Client, error) {
	host, router, err := p2p.NewHost(ctx, priv)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(dstore)
	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, bstore)
	bserv := blockservice.New(bstore, exc)
	dag := merkledag.NewDAGService(bserv)
	resolv := resolver.NewBasicResolver(dag)

	queue, err := queue.NewQueue(ctx, QueueName, dstore)
	if err != nil {
		return nil, err
	}

	prov := simple.NewProvider(ctx, queue, router)
	keys := simple.NewBlockstoreProvider(bstore)
	reprov := simple.NewReprovider(ctx, ReprovideInterval, router, keys)

	system := provider.NewSystem(prov, reprov)
	system.Run()

	for _, info := range dht.GetDefaultBootstrapPeerAddrInfos() {
		go host.Connect(ctx, info)
	}

	if err := p2p.Discovery(ctx, host); err != nil {
		return nil, err
	}

	return &Client{
		DAGService: dag,
		host:       host,
		priv:       priv,
		bstore:     bstore,
		dstore:     dstore,
		resolv:     resolv,
		router:     router,
		system:     system,
	}, nil
}

// PeerID returns the peer id of the client.
func (c *Client) PeerID() (peer.ID, error) {
	return peer.IDFromPrivateKey(c.priv)
}

// ResolvePath resolves the node from the given path.
func (c *Client) ResolvePath(ctx context.Context, p path.Path) (ipld.Node, error) {
	return c.resolv.ResolvePath(ctx, p)
}
