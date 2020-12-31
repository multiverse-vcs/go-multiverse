package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
)

// PullArgs contains the args.
type PullArgs struct {
	// Root is the repo root path.
	Root string
	// Head is the CID of the repo head.
	Head cid.Cid
	// ID is the CID of the commit to pull.
	ID cid.Cid
}

// PullReply contains the reply.
type PullReply struct {
	// ID is the CID of the merged commits.
	ID cid.Cid
}

// Pull merges changes into the repo head.
func (s *Service) Pull(args *PullArgs, reply *PullReply) error {
	ctx := context.Background()

	diffs, err := core.Status(ctx, s.dag, args.Root, args.Head)
	if err != nil {
		return err
	}

	if len(diffs) != 0 {
		return errors.New("repo has uncommitted changes")
	}

	if err := merkledag.FetchGraph(ctx, args.ID, s.dag); err != nil {
		return err
	}

	base, err := core.MergeBase(ctx, s.dag, args.Head, args.ID)
	if err != nil {
		return err
	}

	if base == args.ID {
		return errors.New("local is ahead of remote")
	}

	merge, err := core.Merge(ctx, s.dag, args.Head, base, args.ID)
	if err != nil {
		return err
	}

	reply.ID = merge.Cid()
	return core.Write(ctx, s.dag, args.Root, merge)
}
