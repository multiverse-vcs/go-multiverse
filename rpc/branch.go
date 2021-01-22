package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// BranchArgs contains the args.
type BranchArgs struct {
	// Repo is the CID of the repo.
	Repo cid.Cid
	// Branch is the name of the branch.
	Branch string
	// Head is the CID of the branch head.
	Head cid.Cid
}

// BranchReply contains the reply.
type BranchReply struct {
	// Repo is the CID of the repo.
	Repo cid.Cid
	// Branches is the map of repo branch heads.
	Branches map[string]cid.Cid
}

// ListBranches returns the repo branches.
func (s *Service) ListBranches(args *BranchArgs, reply *BranchReply) error {
	ctx := context.Background()

	repo, err := data.GetRepository(ctx, s.client, args.Repo)
	if err != nil {
		return err
	}

	reply.Repo = args.Repo
	reply.Branches = repo.Branches
	return nil
}

// CreateBranch creates a new branch.
func (s *Service) CreateBranch(args *BranchArgs, reply *BranchReply) error {
	ctx := context.Background()

	repo, err := data.GetRepository(ctx, s.client, args.Repo)
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

	id, err := data.PinRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}
	s.client.Unpin(ctx, args.Repo, true)

	reply.Repo = id
	reply.Branches = repo.Branches
	return nil
}

// DeleteBranch deletes an existing branch.
func (s *Service) DeleteBranch(args *BranchArgs, reply *BranchReply) error {
	ctx := context.Background()

	repo, err := data.GetRepository(ctx, s.client, args.Repo)
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

	id, err := data.PinRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}
	s.client.Unpin(ctx, args.Repo, true)

	reply.Repo = id
	reply.Branches = repo.Branches
	return nil
}
