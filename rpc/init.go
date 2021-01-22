package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// InitArgs contains the args.
type InitArgs struct {
	Name string
}

// InitReply contains the reply.
type InitReply struct {
	// Repo is the CID of the repository.
	Repo cid.Cid
}

// Init creates a new empty repository.
func (s *Service) Init(args *InitArgs, reply *InitReply) error {
	ctx := context.Background()

	if args.Name == "" {
		return errors.New("name cannot be empty")
	}

	repo := data.NewRepository(args.Name)

	id, err := data.PinRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}

	reply.Repo = id
	return nil
}
