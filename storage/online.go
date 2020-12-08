package storage

import (
	"context"

	"github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// Online initializes a p2p host for the underlying blockservice.
func (s *Store) Online(ctx context.Context) error {
	priv, err := s.ReadKey()
	if err != nil {
		return err
	}

	s.Host, s.Router, err = p2p.NewHost(ctx, priv)
	if err != nil {
		return err
	}

	net := bsnet.NewFromIpfsHost(s.Host, s.Router)
	exc := bitswap.New(ctx, net, s.bstore)

	bserv := blockservice.New(s.bstore, exc)
	s.Dag = merkledag.NewDAGService(bserv)

	p2p.Bootstrap(ctx, s.Host)

	if err := p2p.Discovery(ctx, s.Host); err != nil {
		return err
	}

	return s.Router.Bootstrap(ctx)
}
