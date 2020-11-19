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

	commit, ok := node.(*object.Commit)
	if !ok {
		return errors.New("invalid commit")
	}

	link, _, err := commit.ResolveLink([]string{"tree"})
	if err != nil {
		return err
	}

	tree, err := link.GetNode(c.ctx, c.dag)
	if err != nil {
		return err
	}

	return c.Write(c.config.Root, tree)
}
