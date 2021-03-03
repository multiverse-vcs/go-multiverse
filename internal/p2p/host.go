package p2p

import (
	"context"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	noise "github.com/libp2p/go-libp2p-noise"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
)

const (
	// LowWater is the minimum amount of connections to keep.
	LowWater = 100
	// HighWater is the maximum amount of connections to keep.
	HighWater = 400
	// GracePeriod is how long wait to consider a connection active.
	GracePeriod = time.Minute
)

// NewHost returns a new libp2p host and router.
func NewHost(ctx context.Context, priv crypto.PrivKey, listenAddr []string) (host.Host, routing.Routing, error) {
	var router routing.Routing
	var err error

	host, err := libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(listenAddr...),
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
