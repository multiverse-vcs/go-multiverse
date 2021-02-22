package author

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// SearchArgs contains the args.
type SearchArgs struct {
	// PeerID is the peer ID of the author.
	PeerID string `json:"peerID"`
}

// SearchReply contains the reply
type SearchReply struct {
	// Author is the author profile.
	Author *object.Author `json:"author"`
}

// Search returns the author for the given peer ID.
func (s *Service) Search(args *SearchArgs, reply *SearchReply) error {
	ctx := context.Background()

	peerID, err := peer.Decode(args.PeerID)
	if err != nil {
		return err
	}

	authorID, err := s.Namesys.Resolve(ctx, peerID)
	if err != nil {
		return err
	}

	author, err := object.GetAuthor(ctx, s.Peer.DAG, authorID)
	if err != nil {
		return err
	}

	reply.Author = author
	return nil
}
