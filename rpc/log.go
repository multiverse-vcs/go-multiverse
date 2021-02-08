package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// LogArgs contains the args.
type LogArgs struct {
	// Name is the repo name
	Name string
	// Branch is the name of the repo branch.
	Branch string
	// Limit is the number of commits to log.
	Limit int
}

// LogReply contains the reply.
type LogReply struct {
	IDs     []cid.Cid
	Commits []*data.Commit
}

// Log returns the changes between the working directory and repo head.
func (s *Service) Log(args *LogArgs, reply *LogReply) error {
	ctx := context.Background()
	cfg := s.node.Config()

	id, ok := cfg.Author.Repositories[args.Name]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := data.GetRepository(ctx, s.node, id)
	if err != nil {
		return err
	}

	head, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	var ids []cid.Cid
	visit := func(id cid.Cid) bool {
		if args.Limit > -1 && len(ids) >= args.Limit {
			return false
		}

		ids = append(ids, id)
		return true
	}

	if err := core.Walk(ctx, s.node, head, visit); err != nil {
		return err
	}

	var commits []*data.Commit
	for _, id := range ids {
		commit, err := data.GetCommit(ctx, s.node, id)
		if err != nil {
			return err
		}

		commits = append(commits, commit)
	}

	reply.IDs = ids
	reply.Commits = commits
	return nil
}
