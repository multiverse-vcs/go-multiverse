package rpc

import (
	"context"
	"errors"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/multiverse-vcs/go-multiverse/pkg/p2p"
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
)

// PushArgs contains the args.
type PushArgs struct {
	// Remote is the repository path.
	Remote remote.Path
	// Branch is the name of the branch to update.
	Branch string
	// Pack contains nodes to add to the branch.
	Pack []byte
}

// PushReply contains the reply
type PushReply struct{}

// Push updates repository branches and publishes the updated author.
func (s *Service) Push(args *PushArgs, reply *PushReply) error {
	ctx := context.Background()

	key, err := p2p.DecodeKey(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	peerID, err := args.Remote.PeerID()
	if err != nil {
		return err
	}

	name, err := args.Remote.Name()
	if err != nil {
		return err
	}

	if !peerID.MatchesPrivateKey(key) {
		return errors.New("private key does not match")
	}

	author := s.Config.Author
	repoID, ok := author.Repositories[name]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := object.GetRepository(ctx, s.Peer.DAG, repoID)
	if err != nil {
		return err
	}

	old := repo.Branches[args.Branch]
	new, err := remote.LoadPack(ctx, s.Peer.DAG, s.Peer.Blocks, args.Pack, old)
	if err != nil {
		return err
	}

	repo.Branches[args.Branch] = new
	repoID, err = object.AddRepository(ctx, s.Peer.DAG, repo)
	if err != nil {
		return err
	}

	author.Repositories[name] = repoID
	if err := s.Config.Write(); err != nil {
		return err
	}

	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	return s.Namesys.Publish(ctx, key, authorID)
}
