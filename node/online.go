package node

import (
	"context"

	"github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// Online initializes a p2p host for the node.
func (n *Node) Online(ctx context.Context, key crypto.PrivKey) error {
	host, router, err := p2p.NewHost(ctx, key)
	if err != nil {
		return err
	}

	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, n.bstore)
	bserv := blockservice.New(n.bstore, exc)

	n.Dag = merkledag.NewDAGService(bserv)
	n.Host = host
	n.router = router

	// TODO replace with multiverse bootstrap peers
	for _, info := range dht.GetDefaultBootstrapPeerAddrInfos() {
		go host.Connect(ctx, info)
	}

	if err := p2p.Discovery(ctx, host); err != nil {
		return err
	}

	return router.Bootstrap(ctx)
}
