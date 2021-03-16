package author

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/internal/key"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// CreateArgs contains the args.
type CreateArgs struct{}

// CreateReply contains the reply
type CreateReply struct {
	// Peer is the peer ID of the author.
	Peer string `json:"peer"`
}

// Create generates a new private key and publishes an empty author object.
func (s *Service) Create(args *CreateArgs, reply *CreateReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priv, err := key.Generate()
	if err != nil {
		return err
	}

	peerID, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return err
	}

	if err := s.Keystore.Put(peerID.Pretty(), priv); err != nil {
		return err
	}

	author := object.NewAuthor()
	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	reply.Peer = peerID.Pretty()
	return s.Namesys.Publish(ctx, priv, authorID)
}
