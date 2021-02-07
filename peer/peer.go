// Package peer implements a peer client.
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
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// Client manages peer services.
type Client struct {
	ipld.DAGService
	host    host.Host
	config  *Config
	bstore  blockstore.Blockstore
	dstore  datastore.Batching
	namesys routing.ValueStore
	provsys provider.System
	resolv  *resolver.Resolver
	router  routing.Routing
}

// New returns a peer with a p2p host and persistent storage.
func New(ctx context.Context, dstore datastore.Batching, config *Config) (*Client, error) {
	priv, err := p2p.DecodeKey(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	host, router, err := p2p.NewHost(ctx, priv)
	if err != nil {
		return nil, err
	}

	namesys, err := p2p.NewNamesys(ctx, host)
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
	provsys.Run()

	for _, info := range dht.GetDefaultBootstrapPeerAddrInfos() {
		go host.Connect(ctx, info)
	}

	if err := p2p.Discovery(ctx, host); err != nil {
		return nil, err
	}

	return &Client{
		DAGService: dag,
		host:       host,
		config:     config,
		bstore:     bstore,
		dstore:     dstore,
		resolv:     resolv,
		router:     router,
		provsys:    provsys,
		namesys:    namesys,
	}, nil
}

// Authors returns the authors api.
func (c *Client) Authors() *authors {
	return (*authors)(c)
}

// Config returns the peer config.
func (c *Client) Config() *Config {
	return c.config
}

// PeerID returns the peer id of the client.
func (c *Client) PeerID() peer.ID {
	return c.host.ID()
}

// ResolvePath resolves the node from the given path.
func (c *Client) ResolvePath(ctx context.Context, p path.Path) (ipld.Node, error) {
	return c.resolv.ResolvePath(ctx, p)
}
