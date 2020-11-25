package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	//"github.com/whyrusleeping/mdns"
)

var DiscoveryInterval = time.Duration(10) * time.Second

const DiscoveryConnectTimeout = time.Second * 30

type discoveryHandler struct {
	ctx  context.Context
	host host.Host
}

func (dh *discoveryHandler) HandlePeerFound(p peer.AddrInfo) {
	fmt.Println("connecting to discovered peer: ", p)
	ctx, cancel := context.WithTimeout(dh.ctx, DiscoveryConnectTimeout)
	defer cancel()
	if err := dh.host.Connect(ctx, p); err != nil {
		fmt.Printf("failed to connect to peer %s found by discovery: %s\n", p.ID, err)
	}
}

// Discovery starts discovering peers through mdns.
func Discovery(ctx context.Context, h host.Host) error {
	//mdns.DisableLogging = false

	service, err := discovery.NewMdnsService(ctx, h, DiscoveryInterval, discovery.ServiceTag)
	if err != nil {
		return err
	}

	handler := discoveryHandler{ctx, h}
	service.RegisterNotifee(&handler)

	return nil
}
