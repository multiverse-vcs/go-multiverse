package repo

import (
	"context"
	"errors"
	"path"

	"github.com/multiverse-vcs/go-multiverse/internal/p2p"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// CreateArgs contains the args.
type CreateArgs struct {
	// Name is the repository name.
	Name string `json:"name"`
}

// CreateReply contains the reply
type CreateReply struct {
	// Remote is the repository path
	Remote string `json:"remote"`
}

// Create creates a new repository.
func (s *Service) Create(args *CreateArgs, reply *CreateReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if args.Name == "" {
		return errors.New("name cannot be empty")
	}

	key, err := p2p.DecodeKey(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	author := s.Config.Author
	if _, ok := author.Repositories[args.Name]; ok {
		return errors.New("repository already exists")
	}

	repo := object.NewRepository()
	repoID, err := object.AddRepository(ctx, s.Peer.DAG, repo)
	if err != nil {
		return err
	}

	author.Repositories[args.Name] = repoID
	if err := s.Config.Write(); err != nil {
		return err
	}

	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	reply.Remote = path.Join(s.Peer.Host.ID().Pretty(), args.Name)
	return s.Namesys.Publish(ctx, key, authorID)
}
