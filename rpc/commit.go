package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
)

// CommitArgs contains the args.
type CommitArgs struct {
	// Name is the repo name.
	Name string
	// Branch is the name of the branch to update.
	Branch string
	// Root is the repo root path.
	Root string
	// Ignore is a list of paths to ignore.
	Ignore []string
	// Message is the description of changes.
	Message string
	// Parents contains the parent CIDs.
	Parents []cid.Cid
}

// CommitReply contains the reply.
type CommitReply struct {
	// ID is the CID of the commit.
	ID cid.Cid
}

// Commit records changes to the repo
func (s *Service) Commit(args *CommitArgs, reply *CommitReply) error {
	ctx := context.Background()

	if args.Name == "" {
		return errors.New("name cannot be empty")
	}

	if args.Branch == "" {
		return errors.New("branch cannot be empty")
	}

	repo, err := s.node.GetRepositoryOrDefault(ctx, args.Name)
	if err != nil {
		return err
	}

	// TODO verify that parent cids are valid

	id, err := core.Commit(ctx, s.node, args.Root, args.Ignore, args.Message, args.Parents...)
	if err != nil {
		return err
	}

	reply.ID = id
	repo.Branches[args.Branch] = id

	return s.node.PutRepository(ctx, repo)
}
