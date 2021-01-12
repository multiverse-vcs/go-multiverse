package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// CheckoutArgs contains the args.
type CheckoutArgs struct {
	// Root is the repo root path.
	Root string
	// Name is the name of the repo.
	Name string
	// Branch is the name of the repo branch.
	Branch string
	// Index is the CID of the current commit.
	Index cid.Cid
	// Ignore is a list of paths to ignore.
	Ignore []string
	// ID is the CID of the commit to checkout.
	ID cid.Cid
}

// CheckoutReply contains the reply.
type CheckoutReply struct{}

// Checkout copies the tree of an existing commit to the root.
func (s *Service) Checkout(args *CheckoutArgs, reply *CheckoutReply) error {
	ctx := context.Background()

	equal, err := core.Equal(ctx, s.node, args.Root, args.Ignore, args.Index)
	if err != nil {
		return err
	}

	if !equal {
		return errors.New("uncommitted changes")
	}

	repo, err := s.node.GetRepository(ctx, args.Name)
	if err != nil {
		return err
	}

	head, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	child, err := core.IsAncestor(ctx, s.node, head, args.ID)
	if err != nil {
		return err
	}

	if !child {
		return errors.New("commit is not in branch")
	}

	commit, err := data.GetCommit(ctx, s.node, args.ID)
	if err != nil {
		return err
	}

	tree, err := s.node.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	return core.Write(ctx, s.node, args.Root, tree)
}
