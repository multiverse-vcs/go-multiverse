package rpc

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
)

// CommitArgs contains the args.
type CommitArgs struct {
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

	id, err := core.Commit(ctx, s.dag, args.Root, args.Message, args.Parents...)
	if err != nil {
		return err
	}

	reply.ID = id
	return nil
}
