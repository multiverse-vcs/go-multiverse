package author

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// UnfollowArgs contains the args.
type UnfollowArgs struct {
	// PeerID is the peer ID of the author.
	PeerID string `json:"peerID"`
}

// UnfollowReply contains the reply
type UnfollowReply struct{}

// Unfollow returns the author for the given peer ID.
func (s *Service) Unfollow(args *UnfollowArgs, reply *UnfollowReply) error {
	peerID, err := peer.Decode(args.PeerID)
	if err != nil {
		return err
	}

	if _, err := s.Namesys.Unsubscribe(peerID); err != nil {
		return err
	}

	set := make(map[peer.ID]bool)
	for _, id := range s.Config.Author.Following {
		set[id] = true
	}
	delete(set, peerID)

	var list []peer.ID
	for id := range set {
		list = append(list, id)
	}

	s.Config.Author.Following = list
	return s.Config.Write()
}
