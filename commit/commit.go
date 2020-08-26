package commit

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

// Commit contains metadata for the commit.
type Commit struct {
	Message string  `refmt:"message,omitempty"`
	Changes cid.Cid `refmt:"changes,omitempty"`
	Parent  cid.Cid `refmt:"parent,omitempty"`
}

func init() {
	cbornode.RegisterCborType(Commit{})
}

// NewCommit creates commit.
func NewCommit(message string, changes cid.Cid, parent cid.Cid) *Commit {
	return &Commit{Message: message, Changes: changes, Parent: parent}
}

// Get returns the commit with the matching CID.
func Get(ipfs *core.IpfsNode, id cid.Cid) (*Commit, error) {
	dag, err := ipfs.DAG.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	var node Commit
	if err := cbornode.DecodeInto(dag.RawData(), &node); err != nil {
		return nil, err
	}

	return &node, nil
}

// Node returns an ipld node representation of the commit.
func (c *Commit) Node() (format.Node, error) {
	return cbornode.WrapObject(c, multihash.SHA2_256, -1)
}

// String returns a human readable representation of the commit.
func (c *Commit) String() string {
	return fmt.Sprintf("%s\n\t%s\n", c.Changes, c.Message)
}
