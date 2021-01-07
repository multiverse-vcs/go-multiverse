package rpc

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/core"
)

const (
	StatusAdd    = dagutils.Add
	StatusRemove = dagutils.Remove
	StatusMod    = dagutils.Mod
)

// StatusArgs contains the args.
type StatusArgs struct {
	// Root is the repo root path.
	Root string
	// Head is the CID of the repo head.
	Head cid.Cid
	// Ignore is a list of paths to ignore.
	Ignore []string
}

// StatusReply contains the reply.
type StatusReply struct {
	Diffs map[string]dagutils.ChangeType
}

// Status returns the changes between the working directory and repo head.
func (s *Service) Status(args *StatusArgs, reply *StatusReply) error {
	ctx := context.Background()

	dag := &merkledag.ComboService{
		Read:  s.node,
		Write: dagutils.NewMemoryDagService(),
	}

	changes, err := core.Status(ctx, dag, args.Root, args.Ignore, args.Head)
	if err != nil {
		return err
	}

	diffs := make(map[string]dagutils.ChangeType)
	for _, change := range changes {
		if _, ok := diffs[change.Path]; ok {
			diffs[change.Path] = dagutils.Mod
		} else if change.Path != "" {
			diffs[change.Path] = change.Type
		}
	}

	reply.Diffs = diffs
	return nil
}
