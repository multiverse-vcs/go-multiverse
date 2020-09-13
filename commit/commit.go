package commit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-block-format"
	"github.com/multiformats/go-multihash"
	"github.com/yondero/multiverse/ipfs"
)

// Commit contains metadata for the commit.
type Commit struct {
	ID      cid.Cid `json:"-"`
	Message string  `json:"message,omitempty"`
	Changes cid.Cid `json:"changes,omitempty"`
	Parent  cid.Cid `json:"parent,omitempty"`
}

// NewCommit returns a new commit.
func NewCommit(message string, changes cid.Cid, parent cid.Cid) *Commit {
	return &Commit{Message: message, Changes: changes, Parent: parent}
}

// Get returns the commit with the CID from the blockstore.
func Get(ipfs *ipfs.Node, id cid.Cid) (*Commit, error) {
	b, err := ipfs.Blocks.GetBlock(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	c := Commit{ID: id}
	if err := json.Unmarshal(b.RawData(), &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Add persists the commit to the blockstore.
func (c *Commit) Add(ipfs *ipfs.Node) error {
	b, err := c.Block()
	if err != nil {
		return err
	}

	c.ID = b.Cid()
	if err := ipfs.Blocks.AddBlock(b); err != nil {
		return err
	}

	return nil
}

// Block returns a block representation of the Commit.
func (c *Commit) Block() (blocks.Block, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	return blocks.NewBlockWithCid(data, cid.NewCidV1(cid.Raw, hash))
}

// String returns a human readable representation of the Commit.
func (c *Commit) String() string {
	return fmt.Sprintf("commit %s", c.Changes)
}
