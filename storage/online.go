package storage

import (
	"context"

	"github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// Online initializes a p2p host for the underlying blockservice.
func (s *Store) Online(ctx context.Context) error {
	if s.Host != nil {
		return nil
	}

	// TODO save key instead of generating
	priv, _, err := crypto.GenerateKeyPair(p2p.DefaultKeyType, -1)
	if err != nil {
		return err
	}

	host, router, err := p2p.NewHost(ctx, priv)
	if err != nil {
		return err
	}

	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, s.bstore)

	bserv := blockservice.New(s.bstore, exc)
	dag := merkledag.NewDAGService(bserv)

	s.Dag = dag
	s.Host = host

	return nil
}
