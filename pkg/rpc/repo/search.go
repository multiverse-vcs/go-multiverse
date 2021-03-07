package repo

import (
	"context"
	"errors"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// SearchArgs contains the args.
type SearchArgs struct {
	// Remote is the remote path.
	Remote string `json:"remote"`
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

	parts := strings.Split(args.Remote, "/")
	if len(parts) != 2 {
		return errors.New("invalid remote path")
	}

	pname := parts[0]
	rname := parts[1]

	peerID, err := peer.Decode(pname)
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

	repoID, ok := author.Repositories[rname]
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
