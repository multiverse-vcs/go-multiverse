package rpc

import (
	"context"
	"fmt"
	"sort"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/core"
)

// StatusArgs contains the args.
type StatusArgs struct {
	// Root is the repo root path.
	Root string
	// Head is the CID of the repo head.
	Head cid.Cid
}

// StatusReply contains the reply.
type StatusReply struct {
	Diffs []string
}

// Status returns the changes between the working directory and repo head.
func (s *Service) Status(args *StatusArgs, reply *StatusReply) error {
	ctx := context.Background()

	diffs, err := core.Status(ctx, s.node, args.Root, args.Head)
	if err != nil {
		return err
	}

	paths := make([]string, 0)
	for path := range diffs {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, p := range paths {
		switch diffs[p] {
		case dagutils.Add:
			reply.Diffs = append(reply.Diffs, fmt.Sprintf("\tnew file: %s", p))
		case dagutils.Remove:
			reply.Diffs = append(reply.Diffs, fmt.Sprintf("\tdeleted:  %s", p))
		case dagutils.Mod:
			reply.Diffs = append(reply.Diffs, fmt.Sprintf("\tmodified: %s", p))
		}
	}

	return nil
}
