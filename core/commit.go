package core

import (
	"time"

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
		Worktree: node.Cid(),
	}

	if c.config.Head.Defined() {
		commit.Parents = append(commit.Parents, c.config.Head)
	}

	if err := c.dag.Add(c.ctx, &commit); err != nil {
		return nil, err
	}

	c.config.Head = commit.Cid()
	// TODO write config

	return &commit, nil
}
