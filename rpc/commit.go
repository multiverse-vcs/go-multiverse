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
	// Root is the repo root path.
	Root string
	// Message is the description of changes.
	Message string
	// Parents contains the parent CIDs.
	Parents []cid.Cid
}

// CommitReply contains the reply.
type CommitReply struct {
	ID cid.Cid
}

// Commit records changes to the repo
func (s *Service) Commit(args *CommitArgs, reply *CommitReply) error {
	ctx := context.Background()

	if args.Name == "" {
		return errors.New("name cannot be empty")
	}

	id, err := core.Commit(ctx, s.node, args.Root, args.Message, args.Parents...)
	if err != nil {
		return err
	}

	reply.ID = id
	return s.node.Repo().Set(args.Name, id)
}
