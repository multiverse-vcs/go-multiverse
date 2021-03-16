package author

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/internal/key"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// ImportArgs contains the args.
type ImportArgs struct {
	// Data contains the encoded private key.
	Data string `json:"peer"`
}

// ImportReply contains the reply
type ImportReply struct{}

// Import creates an author from an existing private key.
func (s *Service) Import(args *ImportArgs, reply *ImportReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priv, err := key.Decode(args.Data)
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

	return s.Namesys.Publish(ctx, priv, authorID)
}
