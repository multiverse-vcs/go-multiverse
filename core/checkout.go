package core

import (
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/object"
)

// Checkout writes the tree of the commit to the root.
func (c *Context) Checkout(id cid.Cid) error {
	node, err := c.dag.Get(c.ctx, id)
	if err != nil {
		return err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return errors.New("invalid commit")
	}

	tree, err := c.dag.Get(c.ctx, commit.Tree)
	if err != nil {
		return err
	}

	return c.Write(c.fs.Root(), tree)
}
