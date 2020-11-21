package core

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/object"
)

// Walk traverses the commit history starting at the given id.
func (c *Context) Walk(id cid.Cid, cb func(cid.Cid, *object.Commit) bool) (map[string]*object.Commit, error) {
	history := make(map[string]*object.Commit)

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

		node, err := c.dag.Get(c, id)
		if err != nil {
			return false
		}

		commit, err := object.CommitFromCBOR(node.RawData())
		if err != nil {
			return false
		}

		history[id.KeyString()] = commit
		if cb != nil {
			return cb(id, commit)
		}

		return true
	}

	return history, merkledag.Walk(c, getLinks, id, visit)
}
