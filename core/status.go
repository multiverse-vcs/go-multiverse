package core

import (
	"errors"

	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/object"
)

// Status returns a list of changes between the head and worktree.
func (c *Context) Status() ([]*dagutils.Change, error) {
	tree, err := c.Worktree()
	if err != nil {
		return nil, err
	}

	nodeB, err := c.Add(tree)
	if err != nil {
		return nil, err
	}

	if !c.config.Head.Defined() {
		return dagutils.Diff(c.ctx, c.dag, &merkledag.ProtoNode{}, nodeB)
	}

	node, err := c.dag.Get(c.ctx, c.config.Head)
	if err != nil {
		return nil, err
	}

	commit, ok := node.(*object.Commit)
	if !ok {
		return nil, errors.New("invalid commit")
	}

	link, _, err := commit.ResolveLink([]string{"tree"})
	if err != nil {
		return nil, err
	}

	nodeA, err := link.GetNode(c.ctx, c.dag)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(c.ctx, c.dag, nodeA, nodeB)
}
