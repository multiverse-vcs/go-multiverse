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
	// Root is the repo root path.
	Root string
	// Name is the name of the repo.
	Name string
	// Branch is the name of the repo branch.
	Branch string
	// Branches is map of repo branches.
	Branches map[string]cid.Cid
}

// Clone copies a commit tree to the working directory.
func (s *Service) Clone(args *CloneArgs, reply *CloneReply) error {
	ctx := context.Background()

	node, err := s.node.Get(ctx, args.ID)
	if err != nil {
		return err
	}

	repo, err := data.RepositoryFromCBOR(node.RawData())
	if err != nil {
		return err
	}

	if args.Name == "" {
		args.Name = repo.Name
	}

	if args.Branch == "" {
		args.Branch = repo.DefaultBranch()
	}

	id, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	if err := merkledag.FetchGraphWithDepthLimit(ctx, id, args.Limit, s.node); err != nil {
		return err
	}

	path := filepath.Join(args.Cwd, args.Name)
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	reply.Root = path
	reply.Name = repo.Name
	reply.Branch = args.Branch
	reply.Branches = repo.Branches

	return core.Checkout(ctx, s.node, path, id)
}
