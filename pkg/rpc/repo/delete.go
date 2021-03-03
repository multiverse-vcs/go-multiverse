package repo

import (
	"context"
	"errors"

	"github.com/multiverse-vcs/go-multiverse/internal/p2p"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// DeleteArgs contains the args.
type DeleteArgs struct {
	// Name is the repository name.
	Name string `json:"name"`
}

// DeleteReply contains the reply
type DeleteReply struct{}

// Delete delete an existing repository.
func (s *Service) Delete(args *DeleteArgs, reply *DeleteReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key, err := p2p.DecodeKey(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	author := s.Config.Author
	if _, ok := author.Repositories[args.Name]; !ok {
		return errors.New("repository does not exist")
	}

	delete(author.Repositories, args.Name)
	if err := s.Config.Write(); err != nil {
		return err
	}

	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	return s.Namesys.Publish(ctx, key, authorID)
}
