package author

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// FollowArgs contains the args.
type FollowArgs struct {
	// PeerID is the peer ID of the author.
	PeerID peer.ID `json:"peerID"`
}

// FollowReply contains the reply
type FollowReply struct{}

// Follow returns the author for the given peer ID.
func (s *Service) Follow(args *FollowArgs, reply *FollowReply) error {
	if err := args.PeerID.Validate(); err != nil {
		return err
	}

	if err := s.Namesys.Subscribe(args.PeerID); err != nil {
		return err
	}

	set := make(map[peer.ID]bool)
	for _, id := range s.Config.Author.Following {
		set[id] = true
	}
	set[args.PeerID] = true

	var list []peer.ID
	for id := range set {
		list = append(list, id)
	}

	s.Config.Author.Following = list
	return s.Config.Write()
}
