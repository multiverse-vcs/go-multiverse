package repo

import (
	"context"
	"errors"
	"strings"
	"path"

	merkledag "github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/internal/p2p"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// ForkArgs contains the args.
type ForkArgs struct {
	// Remote is the remote path.
	Remote string `json:"remote"`
	// Name is the new repository name.
	Name string `json:"name"`
}

// ForkReply contains the reply.
type ForkReply struct {
	// Remote is the remote path.
	Remote string `json:"remote"`
}

// Fork returns the repository at the given remote path.
func (s *Service) Fork(args *ForkArgs, reply *ForkReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parts := strings.Split(args.Remote, "/")
	if len(parts) != 2 {
		return errors.New("invalid remote path")
	}

	pname := parts[0]
	rname := parts[1]

	author := s.Config.Author
	rename := rname

	key, err := p2p.DecodeKey(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	peerID, err := peer.Decode(pname)
	if err != nil {
		return err
	}

	sourceID, err := s.Namesys.Search(ctx, peerID)
	if err != nil {
		return err
	}

	source, err := object.GetAuthor(ctx, s.Peer.DAG, sourceID)
	if err != nil {
		return err
	}

	repoID, ok := source.Repositories[rname]
	if !ok {
		return errors.New("repository does not exist")
	}

	if args.Name != "" {
		rename = args.Name
	}

	if _, ok := author.Repositories[rename]; ok {
		return errors.New("repository already exists")
	}

	if err := merkledag.FetchGraph(ctx, repoID, s.Peer.DAG); err != nil {
		return err
	}

	author.Repositories[rename] = repoID
	if err := s.Config.Write(); err != nil {
		return err
	}

	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	reply.Remote = path.Join(s.Peer.Host.ID().Pretty(), rename)
	return s.Namesys.Publish(ctx, key, authorID)
}
