package repo

import (
	"bytes"
	"context"
	"errors"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/merge"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// PushArgs contains the args.
type PushArgs struct {
	// Peer is the author peer ID.
	Peer string `json:"key"`
	// Name is the repository name.
	Name string `json:"name"`
	// branch is the branch name.
	Branch string `json:"branch"`
	// Data contains objects to add.
	Data []byte `json:"data"`
}

// PushReply contains the reply.
type PushReply struct{}

func (s *Service) Push(args *PushArgs, reply *PushReply) error {
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

	repoID, ok := author.Repositories[args.Name]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := object.GetRepository(ctx, s.Peer.DAG, repoID)
	if err != nil {
		return err
	}

	prev := repo.Branches[args.Branch]
	next, err := dag.ReadCar(s.Peer.Blocks, bytes.NewReader(args.Data))
	if err != nil {
		return err
	}

	base, err := merge.Base(ctx, s.Peer.DAG, prev, next)
	if err != nil {
		return err
	}

	if base != prev {
		return errors.New("branches are non-divergent")
	}

	repo.Branches[args.Branch] = next
	repoID, err = object.AddRepository(ctx, s.Peer.DAG, repo)
	if err != nil {
		return err
	}

	author.Repositories[args.Name] = repoID
	authorID, err = object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	return s.Namesys.Publish(ctx, priv, authorID)
}
