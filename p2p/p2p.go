// Package p2p contains methods for working with peer-to-peer networks.
package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-libp2p-secio"
	"github.com/libp2p/go-libp2p-tls"
)

// ListenAddresses is a list of addresses to listen on.
var ListenAddresses = []string{
	"/ip4/0.0.0.0/tcp/9000",
	"/ip4/0.0.0.0/udp/9000/quic",
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
	var err error
	var router routing.Routing

	host, err := libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(ListenAddresses...),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(secio.ID, secio.New),
		libp2p.Transport(libp2pquic.NewTransport),
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

// Bootstrap connects to all peers in the default bootstrap list.
func Bootstrap(ctx context.Context, h host.Host) {
	var wg sync.WaitGroup
	for _, info := range dht.GetDefaultBootstrapPeerAddrInfos() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.Connect(ctx, info)
		}()
	}
	wg.Wait()
}
