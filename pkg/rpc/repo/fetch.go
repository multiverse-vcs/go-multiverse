package repo

import (
	"context"
	"errors"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
)

// FetchArgs contains the args.
type FetchArgs struct {
	// Path is the repository path.
	Remote remote.Path `json:"remote"`
}

// FetchReply contains the reply
type FetchReply struct {
	Repository *object.Repository `json:"repository"`
}

// Fetch returns the branches of the repository.
func (s *Service) Fetch(args *FetchArgs, reply *FetchReply) error {
	ctx := context.Background()

	peerID, err := args.Remote.PeerID()
	if err != nil {
		return err
	}

	name, err := args.Remote.Name()
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

	repoID, ok := author.Repositories[name]
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
