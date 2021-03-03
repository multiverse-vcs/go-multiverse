package author

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// SelfArgs contains the args.
type SelfArgs struct{}

// SelfReply contains the reply
type SelfReply struct {
	// Author is the author object.
	Author *object.Author `json:"author"`
	// PeerID is the peer ID of the server.
	PeerID peer.ID `json:"peerID"`
}

// Self returns the server peer's author profile.
func (s *Service) Self(args *SelfArgs, reply *SelfReply) error {
	reply.Author = s.Config.Author
	reply.PeerID = s.Peer.Host.ID()
	return nil
}
