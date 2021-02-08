package peer

import (
	"context"
	"testing"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	path "github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
	bhost "github.com/libp2p/go-libp2p-blankhost"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// Mock implements the peer interface.
type Mock struct {
	dag     ipld.DAGService
	host    host.Host
	config  *Config
	resolv  *resolver.Resolver
	namesys routing.ValueStore
}

// NewMock returns a new mock node.
func NewMock(t *testing.T, ctx context.Context) *Mock {
	net := swarmt.GenSwarm(t, ctx)
	host := bhost.NewBlankHost(net)
	dag := dagutils.NewMemoryDagService()
	resolv := resolver.NewBasicResolver(dag)

	config, err := NewConfig("")
	if err != nil {
		t.Fatal("failed to create peer config")
	}

	namesys, err := p2p.NewNamesys(ctx, host)
	if err != nil {
		t.Fatal("failed to create peer namesys")
	}

	return &Mock{
		dag:     dag,
		host:    host,
		config:  config,
		resolv:  resolv,
		namesys: namesys,
	}
}

// Authors returns the authors api.
func (n *Mock) Authors() *AuthorsAPI {
	return &AuthorsAPI{n}
}

// Config returns the peer config.
func (n *Mock) Config() *Config {
	return n.config
}

// Dag returns the merkledag api.
func (n *Mock) Dag() ipld.DAGService {
	return n.dag
}

// ID returns the peer ID of the node.
func (n *Mock) ID() peer.ID {
	return n.host.ID()
}

// Namesys returns the name system.
func (n *Mock) Namesys() routing.ValueStore {
	return n.namesys
}

// ResolvePath resolves the node from the given path.
func (n *Mock) ResolvePath(ctx context.Context, p path.Path) (ipld.Node, error) {
	return n.resolv.ResolvePath(ctx, p)
}
