package repo

import (
	"context"
	"errors"
	"path"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// CreateArgs contains the args.
type CreateArgs struct {
	// Peer is the author peer ID.
	Peer string `json:"key"`
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

	if _, ok := author.Repositories[args.Name]; ok {
		return errors.New("repository already exists")
	}

	repo := object.NewRepository()
	repoID, err := object.AddRepository(ctx, s.Peer.DAG, repo)
	if err != nil {
		return err
	}

	author.Repositories[args.Name] = repoID
	authorID, err = object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	reply.Remote = path.Join(s.Peer.Host.ID().Pretty(), args.Name)
	return s.Namesys.Publish(ctx, priv, authorID)
}
