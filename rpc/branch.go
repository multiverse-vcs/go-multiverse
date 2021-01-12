package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
)

// BranchArgs contains the args.
type BranchArgs struct {
	// Name is the repo name.
	Name string
	// Branch is the name of the branch.
	Branch string
	// Head is the CID of the repo head.
	Head cid.Cid
}

// BranchReply contains the reply.
type BranchReply struct {
	Branches map[string]cid.Cid
}

// ListBranches returns the repo branches.
func (s *Service) ListBranches(args *BranchArgs, reply *BranchReply) error {
	ctx := context.Background()

	repo, err := s.node.GetRepository(ctx, args.Name)
	if err != nil {
		return err
	}

	reply.Branches = repo.Branches
	return nil
}

// CreateBranch creates a new branch.
func (s *Service) CreateBranch(args *BranchArgs, reply *BranchReply) error {
	ctx := context.Background()

	repo, err := s.node.GetRepository(ctx, args.Name)
	if err != nil {
		return err
	}

	if args.Branch == "" {
		return errors.New("name cannot be empty")
	}

	if _, ok := repo.Branches[args.Branch]; ok {
		return errors.New("branch already exists")
	}

	repo.Branches[args.Branch] = args.Head
	reply.Branches = repo.Branches
	return s.node.PutRepository(ctx, repo)
}

// DeleteBranch deletes an existing branch.
func (s *Service) DeleteBranch(args *BranchArgs, reply *BranchReply) error {
	ctx := context.Background()

	repo, err := s.node.GetRepository(ctx, args.Name)
	if err != nil {
		return err
	}

	if args.Branch == "" {
		return errors.New("name cannot be empty")
	}

	if _, ok := repo.Branches[args.Branch]; !ok {
		return errors.New("branch does not exists")
	}

	delete(repo.Branches, args.Branch)
	reply.Branches = repo.Branches
	return s.node.PutRepository(ctx, repo)
}
