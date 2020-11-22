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

	if !c.Config.Head.Defined() {
		return dagutils.Diff(c, c.Dag, &merkledag.ProtoNode{}, tree)
	}

	node, err := c.Dag.Get(c, c.Config.Head)
	if err != nil {
		return nil, err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return nil, err
	}

	nodeA, err := c.Dag.Get(c, commit.Tree)
	if err != nil {
		return nil, err
	}

	return dagutils.Diff(c, c.Dag, nodeA, tree)
}
