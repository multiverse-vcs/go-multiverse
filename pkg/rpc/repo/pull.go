package repo

import (
	"bytes"
	"context"
	"errors"
	"strings"

	cid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// PullArgs contains the args.
type PullArgs struct {
	// Remote is the remote path.
	Remote string
	// Branch is the branch name.
	Branch string
	// Refs is a list of known references.
	Refs []cid.Cid
}

// PullReply contains the reply.
type PullReply struct {
	// Data contains objects to add.
	Data []byte
}

func (s *Service) Pull(args *PullArgs, reply *PullReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parts := strings.Split(args.Remote, "/")
	if len(parts) != 2 {
		return errors.New("invalid remote")
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

	head, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	refs := cid.NewSet()
	for _, id := range args.Refs {
		refs.Add(id)
	}

	var data bytes.Buffer
	if err := dag.WriteCar(ctx, s.Peer.DAG, head, refs, &data); err != nil {
		return err
	}

	reply.Data = data.Bytes()
	return nil
}
