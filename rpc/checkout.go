package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
)

// CheckoutArgs contains the args.
type CheckoutArgs struct {
	// Root is the repo root path.
	Root string
	// Head is the CID of the repo head.
	Head cid.Cid
	// ID is the CID of the commit.
	ID cid.Cid
}

// CheckoutReply contains the reply.
type CheckoutReply struct{}

// Checkout copies a commit tree to the working directory.
func (s *Service) Checkout(args *CheckoutArgs, reply *CheckoutReply) error {
	ctx := context.Background()

	diffs, err := core.Status(ctx, s.node, args.Root, args.Head)
	if err != nil {
		return err
	}

	if len(diffs) != 0 {
		return errors.New("repo has uncommitted changes")
	}

	return core.Checkout(ctx, s.node, args.Root, args.ID)
}
