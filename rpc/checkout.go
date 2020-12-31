package rpc

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
)

// CheckoutArgs contains the args.
type CheckoutArgs struct {
	// Root is the repo root path.
	Root string
	// CID is the CID of the commit.
	ID cid.Cid
}

// CheckoutReply contains the reply.
type CheckoutReply struct{}

// Checkout copies a commit tree to the working directory.
func (s *Service) Checkout(args *CheckoutArgs, reply *CheckoutReply) error {
	return core.Checkout(context.Background(), s.dag, args.Root, args.ID)
}
