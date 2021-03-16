package repo

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// SearchArgs contains the args.
type SearchArgs struct {
	// Peer is the author peer ID.
	Peer string `json:"key"`
	// Name is the repository name.
	Name string `json:"name"`
}

// SearchReply contains the reply.
type SearchReply struct {
	// Repository contains repo info.
	Repository *object.Repository `json:"repository"`
}

// Search returns the repository at the given remote path.
func (s *Service) Search(args *SearchArgs, reply *SearchReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	peerID, err := peer.Decode(args.Peer)
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

	repoID, ok := author.Repositories[args.Name]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := object.GetRepository(ctx, s.Peer.DAG, repoID)
	if err != nil {
		return err
	}

	reply.Repository = repo
	return nil
}
