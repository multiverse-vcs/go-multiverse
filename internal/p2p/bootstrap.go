package p2p

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type bootstrapper struct {
	wg sync.WaitGroup
}

// Connect initiates a connection to a peer.
func (bs *bootstrapper) Connect(ctx context.Context, host host.Host, info peer.AddrInfo) {
	defer bs.wg.Done()
	host.Connect(ctx, info)
}

// Bootstrap initiates connections to a list of known peers.
func Bootstrap(ctx context.Context, host host.Host, router routing.Routing) error {
	bs := &bootstrapper{}
	for _, info := range dht.GetDefaultBootstrapPeerAddrInfos() {
		bs.wg.Add(1)
		go bs.Connect(ctx, host, info)
	}
	bs.wg.Wait()
	return router.Bootstrap(ctx)
}
