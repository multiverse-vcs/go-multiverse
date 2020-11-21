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

	if !c.cfg.Head.Defined() {
		return dagutils.Diff(c, c.dag, &merkledag.ProtoNode{}, tree)
	}

	node, err := c.dag.Get(c, c.cfg.Head)
	if err != nil {
		return nil, err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return nil, err
	}

	nodeA, err := c.dag.Get(c, commit.Tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(c, c.dag, nodeA, tree)
}
