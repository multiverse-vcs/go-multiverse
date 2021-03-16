package author

import (
	"github.com/multiverse-vcs/go-multiverse/internal/key"
)

// ExportArgs contains the args.
type ExportArgs struct {
	// Peer is the peer ID of the author.
	Peer string `json:"peer"`
}

// ExportReply contains the reply
type ExportReply struct {
	// Data contains the encoded private key.
	Data string `json:"data"`
}

// Export returns the encoded author private key to use when importing.
func (s *Service) Export(args *ExportArgs, reply *ExportReply) error {
	priv, err := s.Keystore.Get(args.Peer)
	if err != nil {
		return err
	}

	data, err := key.Encode(priv)
	if err != nil {
		return err
	}

	reply.Data = data
	return nil
}
