package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p-noise"
	"github.com/libp2p/go-libp2p-tls"
)

// ListenAddresses is a list of addresses to listen on.
var ListenAddresses = []string{
	"/ip4/0.0.0.0/tcp/9000",
}

const (
	// LowWater is the minimum amount of connections to keep.
	LowWater = 100
	// HighWater is the maximum amount of connections to keep.
	HighWater = 400
	// GracePeriod is how long wait to consider a connection active.
	GracePeriod = time.Minute
)

// NewHost returns a new libp2p host and router.
func NewHost(ctx context.Context, priv crypto.PrivKey) (host.Host, routing.Routing, error) {
	var router routing.Routing
	var err error

	host, err := libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(ListenAddresses...),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
		libp2p.EnableAutoRelay(),
		libp2p.ConnectionManager(connmgr.NewConnManager(
			LowWater,
			HighWater,
			GracePeriod,
		)),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			router, err = dual.New(ctx, h)
			return router, err
		}),
	)

	return host, router, err
}

// Bootstrap initiates connections to a list of known peers.
func Bootstrap(ctx context.Context, host host.Host) {
	var wg sync.WaitGroup
	for _, info := range dht.GetDefaultBootstrapPeerAddrInfos() {
		wg.Add(1)
		go func(info peer.AddrInfo) {
			defer wg.Done()
			host.Connect(ctx, info)
		}(info)
	}
	wg.Wait()
}
