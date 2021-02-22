package author

// SelfArgs contains the args.
type SelfArgs struct {}

// SelfReply contains the reply
type SelfReply struct {
	// PeerID is the peer ID of the server.
	PeerID string `json:"peerID"`
}

// Self returns the server peer's author profile.
func (s *Service) Self(args *SelfArgs, reply *SelfReply) error {
	reply.PeerID = s.Peer.Host.ID().Pretty()
	return nil
}
