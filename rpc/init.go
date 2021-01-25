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
}

// InitReply contains the reply.
type InitReply struct {}

// Init creates a new empty repository.
func (s *Service) Init(args *InitArgs, reply *InitReply) error {
	ctx := context.Background()

	if args.Name == "" {
		return errors.New("name cannot be empty")
	}

	if _, err := s.store.GetCid(args.Name); err == nil {
		return errors.New("repo with name already exists")
	}

	repo := data.NewRepository(args.Name)

	id, err := data.AddRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}

	return s.store.PutCid(repo.Name, id)
}
