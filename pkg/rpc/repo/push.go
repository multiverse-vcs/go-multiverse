package repo

import (
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/internal/p2p"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/merge"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// PushArgs contains the args.
type PushArgs struct {
	// Remote is the remote path.
	Remote string
	// branch is the branch name.
	Branch string
	// Data contains objects to add.
	Data []byte
}

// PushReply contains the reply.
type PushReply struct{}

func (s *Service) Push(args *PushArgs, reply *PushReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parts := strings.Split(args.Remote, "/")
	if len(parts) != 2 {
		return errors.New("invalid remote")
	}

	pname := parts[0]
	rname := parts[1]

	key, err := p2p.DecodeKey(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	peerID, err := peer.Decode(pname)
	if err != nil {
		return err
	}

	if !peerID.MatchesPrivateKey(key) {
		return errors.New("private key does not match")
	}

	author := s.Config.Author
	repoID, ok := author.Repositories[rname]
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

	author.Repositories[rname] = repoID
	if err := s.Config.Write(); err != nil {
		return err
	}

	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	return s.Namesys.Publish(ctx, key, authorID)
}
