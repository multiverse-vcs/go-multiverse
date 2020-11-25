package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

const (
	// DiscoveryInterval is the time between mdns queries.
	DiscoveryInterval = time.Duration(15) * time.Second
	// DiscoveryConnectTimeout is the time to wait for a connection to be made.
	DiscoveryConnectTimeout = time.Second * 30
	// DiscoveryServiceTag is used to identify a group of nodes.
	DiscoveryServiceTag = discovery.ServiceTag
)

type discoveryHandler struct {
	ctx  context.Context
	host host.Host
}

// HandlePeerFound is called when a peer is found through mdns.
func (dh *discoveryHandler) HandlePeerFound(p peer.AddrInfo) {
	ctx, cancel := context.WithTimeout(dh.ctx, DiscoveryConnectTimeout)
	defer cancel()
	dh.host.Connect(ctx, p)
}

// Discovery starts discovering peers through mdns.
func Discovery(ctx context.Context, h host.Host) error {
	service, err := discovery.NewMdnsService(ctx, h, DiscoveryInterval, DiscoveryServiceTag)
	if err != nil {
		return err
	}

	handler := discoveryHandler{ctx, h}
	service.RegisterNotifee(&handler)

	return nil
}
