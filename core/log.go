package core

import (
	"context"
	"errors"
	"io"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/object"
)

// Log prints commit history starting at the current head.
func (c *Context) Log(w io.Writer) error {
	if !c.config.Head.Defined() {
		return nil
	}

	getLinks := func(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
		node, err := c.dag.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		commit, err := object.CommitFromCBOR(node.RawData())
		if err != nil {
			return nil, errors.New("invalid commit")
		}

		return commit.ParentLinks(), nil
	}

	visit := func(id cid.Cid) bool {
		node, err := c.dag.Get(c.ctx, id)
		if err != nil {
			return false
		}

		commit, err := object.CommitFromCBOR(node.RawData())
		if err != nil {
			return false
		}

		commit.Log(w, id, c.config.Head, c.config.Base)
		return true
	}

	return merkledag.Walk(c.ctx, getLinks, c.config.Head, visit)
}
