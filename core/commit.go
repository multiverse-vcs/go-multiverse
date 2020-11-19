package core

import (
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"github.com/multiverse-vcs/go-multiverse/object"
)

// Commit creates a new commit.
func (c *Context) Commit(message string) (cid.Cid, error) {
	tree, err := c.Worktree()
	if err != nil {
		return cid.Cid{}, err
	}

	commit := object.Commit{
		Date:    time.Now(),
		Message: message,
		Tree:    tree.Cid(),
	}

	if c.config.Head.Defined() {
		commit.Parents = append(commit.Parents, c.config.Head)
	}

	node, err := cbornode.WrapObject(&commit, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := c.dag.Add(c.ctx, node); err != nil {
		return cid.Cid{}, err
	}

	c.config.Head = node.Cid()
	// TODO write config

	return node.Cid(), nil
}
