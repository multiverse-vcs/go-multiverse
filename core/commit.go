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
	if c.Config.Base != c.Config.Head {
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

	if c.Config.Base.Defined() {
		commit.Parents = append(commit.Parents, c.Config.Base)
	}

	node, err := cbornode.WrapObject(&commit, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := c.Dag.Add(c, node); err != nil {
		return cid.Cid{}, err
	}

	c.Config.Base = node.Cid()
	c.Config.Head = node.Cid()

	return node.Cid(), nil
}
