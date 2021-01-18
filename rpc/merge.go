package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// MergeArgs contains the args.
type MergeArgs struct {
	// Name is the name of the repo.
	Name string
	// Branch is the name of the repo branch.
	Branch string
	// Root is the repo root path.
	Root string
	// Index is the CID of the current commit.
	Index cid.Cid
	// Ignore is a list of paths to ignore.
	Ignore []string
	// ID is the CID of the commit to merge.
	ID cid.Cid
}

// MergeReply contains the reply.
type MergeReply struct {
	// ID is the CID of the merged commits.
	ID cid.Cid
}

// Merge merges changes into the repo head.
func (s *Service) Merge(args *MergeArgs, reply *MergeReply) error {
	ctx := context.Background()

	equal, err := core.Equal(ctx, s.node, args.Root, args.Ignore, args.Index)
	if err != nil {
		return err
	}

	if !equal {
		return errors.New("uncommitted changes")
	}

	rid, err := s.store.GetCid(args.Name)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.node, rid)
	if err != nil {
		return err
	}

	head, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	if head != args.Index {
		return errors.New("index is behind branch head")
	}

	base, err := core.MergeBase(ctx, s.node, head, args.ID)
	if err != nil {
		return err
	}

	if base == args.ID {
		return errors.New("local is ahead of remote")
	}

	merge, err := core.Merge(ctx, s.node, head, base, args.ID)
	if err != nil {
		return err
	}

	commit := data.NewCommit(merge.Cid(), "merge", head, args.ID)
	id, err := data.AddCommit(ctx, s.node, commit)
	if err != nil {
		return err
	}

	repo.Branches[args.Branch] = id
	reply.ID = id

	rid, err = data.AddRepository(ctx, s.node, repo)
	if err != nil {
		return err
	}

	if err := s.store.PutCid(repo.Name, rid); err != nil {
		return err
	}

	return core.Write(ctx, s.node, args.Root, merge)
}
