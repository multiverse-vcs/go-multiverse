package data

import (
	"encoding/json"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
)

// Commit contains info about changes to a repo.
type Commit struct {
	// Date is the timestamp of when the commit was created.
	Date time.Time `json:"date"`
	// Message is a description of the changes.
	Message string `json:"message"`
	// Parents is a list of the parent commit CIDs.
	Parents []cid.Cid `json:"parents"`
	// Tree is the root CID of the repo file tree.
	Tree cid.Cid `json:"tree"`
	// Metadata contains additional data.
	Metadata map[string]string `json:"metadata"`
}

// CommitFromJON decodes a commit from json.
func CommitFromJSON(data []byte) (*Commit, error) {
	var commit Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// CommitFromNode decodes a commit from an ipld node.
func CommitFromCBOR(data []byte) (*Commit, error) {
	var commit Commit
	if err := cbornode.DecodeInto(data, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// ParentLinks returns parent ipld links.
func (c *Commit) ParentLinks() []*ipld.Link {
	out := make([]*ipld.Link, 0)
	for _, p := range c.Parents {
		out = append(out, &ipld.Link{Cid: p})
	}

	return out
}
