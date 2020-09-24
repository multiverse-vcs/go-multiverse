package commit

import (
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

// Commit contains repo changes.
type Commit struct {
	// Message about the changes in the commit.
	Message string    `json:"message"`
	// Commit contents tree.
	Tree cid.Cid      `json:"tree"`
	// Merkle root of the parent commit.
	Parents []cid.Cid `json:"parents"`
}

func init() {
	cbornode.RegisterCborType(Commit{})
}

// FromNode returns a commit from the node.
func FromNode(node format.Node) (*Commit, error) {
	var c Commit
	if err := cbornode.DecodeInto(node.RawData(), &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Node returns a node from the commit.
func (c *Commit) Node() (format.Node, error) {
	return cbornode.WrapObject(c, multihash.SHA2_256, -1)
}