package commit

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

// Signature contains info about who created the commit.
type Signature struct {
	// Name is the name of the person who created the commit.
	Name string `json:"name"`
	// Email is an address that can be used to contact the committer.
	Email string `json:"email"`
	// When is the timestamp of when the commit was created.
	When time.Time `json:"when"`
}

// Commit contains info about changes to a repo.
type Commit struct {
	// Author is the person that created the commit.
	Author Signature `json:"author"`
	// Committer is the person that performed the commit.
	Committer Signature `json:"committer"`
	// Message is a description of the changes.
	Message string `json:"message"`
	// Tree is the current state of the repo files.
	Tree cid.Cid `json:"tree"`
	// Parents contains the CIDs of parent commits.
	Parents []cid.Cid `json:"parents"`
}

// FromNode returns a commit from the node.
func FromNode(node format.Node) (*Commit, error) {
	b, err := node.(*cbornode.Node).MarshalJSON()
	if err != nil {
		return nil, err
	}

	var c Commit
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Node returns a node from the commit.
func (c *Commit) Node() (format.Node, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return cbornode.FromJSON(bytes.NewReader(b), multihash.SHA2_256, -1)
}
