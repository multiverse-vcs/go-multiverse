// Package peer implements a peer client.
package peer

import (
	"context"

	ipld "github.com/ipfs/go-ipld-format"
	path "github.com/ipfs/go-path"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
)

// Peer is used to access p2p network resources.
type Peer interface {
	// Authors returns the authors api.
	Authors() *AuthorsAPI
	// Config returns the peer config.
	Config() *Config
	// Dag returns the merkledag api.
	Dag() ipld.DAGService
	// ID returns the peer ID of the node.
	ID() peer.ID
	// Namesys returns the name system.
	Namesys() routing.ValueStore
	// ResolvePath resolves the node from the given path.
	ResolvePath(context.Context, path.Path) (ipld.Node, error)
}
