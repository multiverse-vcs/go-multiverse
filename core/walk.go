package core

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// WalkFun is called for each commit visited by walk.
type WalkFun func(cid.Cid, *data.Commit) bool

// Walk traverses the commit history starting at the given id.
func Walk(ctx context.Context, dag ipld.DAGService, id cid.Cid, cb WalkFun) (map[string]*data.Commit, error) {
	history := make(map[string]*data.Commit)

	// perform a depth first traversal of parent commits
	getLinks := func(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
		commit, ok := history[id.KeyString()]
		if !ok {
			return nil, errors.New("commit not found")
		}

		return commit.ParentLinks(), nil
	}

	visit := func(id cid.Cid) bool {
		if _, ok := history[id.KeyString()]; ok {
			return false
		}

		node, err := dag.Get(ctx, id)
		if err != nil {
			return false
		}

		commit, err := data.CommitFromCBOR(node.RawData())
		if err != nil {
			return false
		}

		history[id.KeyString()] = commit
		if cb != nil {
			return cb(id, commit)
		}

		return true
	}

	return history, merkledag.Walk(ctx, getLinks, id, visit)
}
