package core

import (
	"time"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/object"
)

// Commit creates a new commit.
func (c *Context) Commit(message string) (*object.Commit, error) {
	tree, err := c.Worktree()
	if err != nil {
		return nil, err
	}

	node, err := c.Add(tree)
	if err != nil {
		return nil, err
	}

	commit := object.Commit{
		Date:     time.Now(),
		Message:  message,
		Parents:  []cid.Cid{c.config.Head},
		Worktree: node.Cid(),
	}

	if err := c.dag.Add(c.ctx, &commit); err != nil {
		return nil, err
	}

	c.config.Head = commit.Cid()

	return &commit, nil
}
