package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// PullArgs contains the args.
type PullArgs struct {
	// Root is the repo root path.
	Root string
	// Head is the CID of the repo head.
	Head cid.Cid
	// Ignore is a list of paths to ignore.
	Ignore []string
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
	dag := dagutils.NewMemoryDagService()

	tree, err := core.Add(ctx, dag, args.Root, args.Ignore)
	if err != nil {
		return err
	}

	node, err := s.node.Get(ctx, args.Head)
	if err != nil {
		return err
	}

	commit, err := data.CommitFromCBOR(node.RawData())
	if err != nil {
		return err
	}

	if tree.Cid() != commit.Tree {
		return errors.New("uncommitted changes")
	}

	// node, err := s.node.Get(ctx, args.ID)
	// if err != nil {
	// 	return err
	// }

	// _, err := data.CommitFromCBOR(node.RawData())
	// if err != nil {
	// 	return err
	// }

	if err := merkledag.FetchGraph(ctx, args.ID, s.node); err != nil {
		return err
	}

	base, err := core.MergeBase(ctx, s.node, args.Head, args.ID)
	if err != nil {
		return err
	}

	if base == args.ID {
		return errors.New("local is ahead of remote")
	}

	merge, err := core.Merge(ctx, s.node, args.Head, base, args.ID)
	if err != nil {
		return err
	}

	reply.ID = merge.Cid()
	return core.Write(ctx, s.node, args.Root, merge)
}
