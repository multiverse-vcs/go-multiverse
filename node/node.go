package node

import (
	"context"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-provider"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
)

// Prefix is the datastore key prefix.
var Prefix = datastore.NewKey("multiverse")

// Node manages peer services.
type Node struct {
	ipld.DAGService
	host   host.Host
	bstore blockstore.Blockstore
	dstore datastore.Batching
	resolv *resolver.Resolver
	router routing.Routing
	system provider.System
}

// ResolvePath resolves the node from the given path.
func (n *Node) ResolvePath(ctx context.Context, p path.Path) (ipld.Node, error) {
	return n.resolv.ResolvePath(ctx, p)
}

// Repo returns a repo API.
func (n *Node) Repo() *repo {
	return &repo{namespace.Wrap(n.dstore, Prefix)}
}
