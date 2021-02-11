package peer

import (
	"context"

	bitswap "github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	blockservice "github.com/ipfs/go-blockservice"
	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	provider "github.com/ipfs/go-ipfs-provider"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	path "github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	namesys "github.com/libp2p/go-libp2p-pubsub-router"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

var _ Peer = (*Node)(nil)

// Node implements the peer interface
type Node struct {
	dag     ipld.DAGService
	host    host.Host
	config  *Config
	bstore  blockstore.Blockstore
	dstore  datastore.Batching
	namesys *namesys.PubsubValueStore
	provsys provider.System
	resolv  *resolver.Resolver
	router  routing.Routing
}

// New returns a peer with a p2p host and persistent storage.
func New(ctx context.Context, dstore datastore.Batching, config *Config) (*Node, error) {
	priv, err := p2p.DecodeKey(config.PrivateKey)
	if err != nil {
		return nil, err
	}

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

	provsys, err := p2p.NewProvider(ctx, dstore, bstore, router)
	if err != nil {
		return nil, err
	}

	p2p.Bootstrap(ctx, host)
	provsys.Run()

	namesys, err := p2p.NewNamesys(ctx, host, router)
	if err != nil {
		return nil, err
	}

	if err := p2p.Discovery(ctx, host); err != nil {
		return nil, err
	}

	if err := router.Bootstrap(ctx); err != nil {
		return nil, err
	}

	return &Node{
		dag:     dag,
		host:    host,
		config:  config,
		bstore:  bstore,
		dstore:  dstore,
		resolv:  resolv,
		router:  router,
		provsys: provsys,
		namesys: namesys,
	}, nil
}

// Authors returns the authors api.
func (n *Node) Authors() *AuthorsAPI {
	return &AuthorsAPI{n}
}

// Config returns the peer config.
func (n *Node) Config() *Config {
	return n.config
}

// Connect connects to the peer with the given ID.
func (n *Node) Connect(ctx context.Context, id peer.ID) error {
	info, err := n.router.FindPeer(ctx, id)
	if err != nil {
		return nil
	}

	return n.host.Connect(ctx, info)
}

// Dag returns the merkledag api.
func (n *Node) Dag() ipld.DAGService {
	return n.dag
}

// ID returns the peer ID of the node.
func (n *Node) ID() peer.ID {
	return n.host.ID()
}

// Namesys returns the name system.
func (n *Node) Namesys() *namesys.PubsubValueStore {
	return n.namesys
}

// ResolvePath resolves the node from the given path.
func (n *Node) ResolvePath(ctx context.Context, p path.Path) (ipld.Node, error) {
	return n.resolv.ResolvePath(ctx, p)
}
