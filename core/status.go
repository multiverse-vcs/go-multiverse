package core

import (
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

	if !c.config.Head.Defined() {
		return dagutils.Diff(c.ctx, c.dag, &merkledag.ProtoNode{}, tree)
	}

	node, err := c.dag.Get(c.ctx, c.config.Head)
	if err != nil {
		return nil, err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return nil, err
	}

	nodeA, err := c.dag.Get(c.ctx, commit.Tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(c.ctx, c.dag, nodeA, tree)
}
