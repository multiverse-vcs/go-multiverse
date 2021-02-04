package rpc

import (
	"context"
	"errors"

	"github.com/multiverse-vcs/go-multiverse/data"
)

// InitArgs contains the args.
type InitArgs struct {
	// Name is the name of the repo.
	Name string
	// Branch is the name of the default branch.
	Branch string
}

// InitReply contains the reply.
type InitReply struct{}

// Init creates a new empty repository.
func (s *Service) Init(args *InitArgs, reply *InitReply) error {
	ctx := context.Background()
	cfg := s.client.Config()

	if args.Name == "" {
		return errors.New("name cannot be empty")
	}

	if args.Branch == "" {
		return errors.New("branch cannot be empty")
	}

	if _, ok := cfg.Author.Repositories[args.Name]; ok {
		return errors.New("repo with name already exists")
	}

	repo := data.NewRepository(args.Name)
	repo.DefaultBranch = args.Branch

	id, err := data.AddRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}
	cfg.Author.Repositories[args.Name] = id

	return cfg.Save()
}
