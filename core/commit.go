package core

import (
	"errors"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"github.com/multiverse-vcs/go-multiverse/object"
)

// Commit creates a new commit.
func (c *Context) Commit(message string) (cid.Cid, error) {
	if c.cfg.Base != c.cfg.Head {
		return cid.Cid{}, errors.New("base is behind head")
	}

	tree, err := c.Worktree()
	if err != nil {
		return cid.Cid{}, err
	}

	commit := object.Commit{
		Date:    time.Now(),
		Message: message,
		Tree:    tree.Cid(),
	}

	if c.cfg.Base.Defined() {
		commit.Parents = append(commit.Parents, c.cfg.Base)
	}

	node, err := cbornode.WrapObject(&commit, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := c.dag.Add(c, node); err != nil {
		return cid.Cid{}, err
	}

	c.cfg.Base = node.Cid()
	c.cfg.Head = node.Cid()

	return node.Cid(), nil
}
