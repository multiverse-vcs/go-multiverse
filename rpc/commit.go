package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
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
	// Parent is the CID of the parent commit.
	Parent cid.Cid
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

	head, ok := repo.Branches[args.Branch]
	if ok && args.Parent != head {
		return errors.New("branch is ahead of parent")
	}

	tree, err := core.Add(ctx, s.node, args.Root, args.Ignore)
	if err != nil {
		return err
	}

	commit := data.NewCommit(tree.Cid(), args.Message)
	if args.Parent.Defined() {
		commit.Parents = append(commit.Parents, args.Parent)
	}

	id, err := data.AddCommit(ctx, s.node, commit)
	if err != nil {
		return err
	}

	reply.ID = id
	repo.Branches[args.Branch] = id

	return s.node.PutRepository(ctx, repo)
}
