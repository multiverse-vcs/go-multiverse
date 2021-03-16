package repo

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// DeleteArgs contains the args.
type DeleteArgs struct {
	// Peer is the author peer ID.
	Peer string `json:"key"`
	// Name is the repository name.
	Name string `json:"name"`
}

// DeleteReply contains the reply
type DeleteReply struct{}

// Delete delete an existing repository.
func (s *Service) Delete(args *DeleteArgs, reply *DeleteReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priv, err := s.Keystore.Get(args.Peer)
	if err != nil {
		return err
	}

	peerID, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return err
	}

	authorID, err := s.Namesys.Search(ctx, peerID)
	if err != nil {
		return err
	}

	author, err := object.GetAuthor(ctx, s.Peer.DAG, authorID)
	if err != nil {
		return err
	}

	if _, ok := author.Repositories[args.Name]; !ok {
		return errors.New("repository does not exist")
	}

	delete(author.Repositories, args.Name)
	authorID, err = object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	return s.Namesys.Publish(ctx, priv, authorID)
}
