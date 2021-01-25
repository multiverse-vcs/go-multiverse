package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

// CheckoutArgs contains the args.
type CheckoutArgs struct {
	// Name is the name of the repo.
	Name string
	// Root is the repo root path.
	Root string
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

	equal, err := core.Equal(ctx, s.client, args.Root, args.Ignore, args.Index)
	if err != nil {
		return err
	}

	if !equal {
		return errors.New("uncommitted changes")
	}

	id, err := s.store.GetCid(args.Name)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return err
	}

	head, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	child, err := core.IsAncestor(ctx, s.client, head, args.ID)
	if err != nil {
		return err
	}

	if !child {
		return errors.New("commit is not in branch")
	}

	commit, err := data.GetCommit(ctx, s.client, args.ID)
	if err != nil {
		return err
	}

	tree, err := s.client.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	return unixfs.Write(ctx, s.client, args.Root, tree)
}
