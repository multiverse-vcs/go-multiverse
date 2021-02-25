package author

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// UnfollowArgs contains the args.
type UnfollowArgs struct {
	// PeerID is the peer ID of the author.
	PeerID peer.ID `json:"peerID"`
}

// UnfollowReply contains the reply
type UnfollowReply struct{}

// Unfollow returns the author for the given peer ID.
func (s *Service) Unfollow(args *UnfollowArgs, reply *UnfollowReply) error {
	if err := args.PeerID.Validate(); err != nil {
		return err
	}

	if _, err := s.Namesys.Unsubscribe(args.PeerID); err != nil {
		return err
	}

	set := make(map[peer.ID]bool)
	for _, id := range s.Config.Author.Following {
		set[id] = true
	}
	delete(set, args.PeerID)

	var list []peer.ID
	for id := range set {
		list = append(list, id)
	}

	s.Config.Author.Following = list
	return s.Config.Write()
}
