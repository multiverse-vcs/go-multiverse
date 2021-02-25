package author

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// FollowArgs contains the args.
type FollowArgs struct {
	// PeerID is the peer ID of the author.
	PeerID string `json:"peerID"`
}

// FollowReply contains the reply
type FollowReply struct{}

// Follow returns the author for the given peer ID.
func (s *Service) Follow(args *FollowArgs, reply *FollowReply) error {
	peerID, err := peer.Decode(args.PeerID)
	if err != nil {
		return err
	}

	if err := s.Namesys.Subscribe(peerID); err != nil {
		return err
	}

	set := make(map[peer.ID]bool)
	for _, id := range s.Config.Author.Following {
		set[id] = true
	}
	set[peerID] = true

	var list []peer.ID
	for id := range set {
		list = append(list, id)
	}

	s.Config.Author.Following = list
	return s.Config.Write()
}
