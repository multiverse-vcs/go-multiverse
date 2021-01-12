package rpc

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// CloneArgs contains the args.
type CloneArgs struct {
	// Cwd is the current working directory.
	Cwd string
	// ID is the CID of the repo.
	ID cid.Cid
	// Limit is the number of children to fetch.
	Limit int
	// Name is the name of the directory to create.
	Name string
	// Branch is the name of the branch to clone.
	Branch string
}

// CloneReply contains the reply.
type CloneReply struct {
	// ID is the CID of the commit.
	ID cid.Cid
	// Root is the repo root path.
	Root string
	// Name is the name of the repo.
	Name string
	// Branch is the name of the repo branch.
	Branch string
}

// Clone copies a commit tree to the working directory.
func (s *Service) Clone(args *CloneArgs, reply *CloneReply) error {
	ctx := context.Background()

	repo, err := data.GetRepository(ctx, s.node, args.ID)
	if err != nil {
		return err
	}

	if args.Name == "" {
		args.Name = repo.Name
	}

	id, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	if err := merkledag.FetchGraphWithDepthLimit(ctx, id, args.Limit, s.node); err != nil {
		return err
	}

	commit, err := data.GetCommit(ctx, s.node, id)
	if err != nil {
		return err
	}

	tree, err := s.node.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	path := filepath.Join(args.Cwd, args.Name)
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	reply.ID = id
	reply.Root = path
	reply.Name = repo.Name
	reply.Branch = args.Branch

	return core.Write(ctx, s.node, path, tree)
}
