package object

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
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

// GetCommit returns the commit with the given CID.
func GetCommit(ctx context.Context, ds ipld.NodeGetter, id cid.Cid) (*Commit, error) {
	node, err := ds.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return CommitFromCBOR(node.RawData())
}

// GetCommitTree returns the tree of the commit with the given CID.
func GetCommitTree(ctx context.Context, ds ipld.NodeGetter, id cid.Cid) (ipld.Node, error) {
	commit, err := GetCommit(ctx, ds, id)
	if err != nil {
		return nil, err
	}

	return ds.Get(ctx, commit.Tree)
}

// AddCommit adds a commit to the given dag.
func AddCommit(ctx context.Context, ds ipld.NodeAdder, commit *Commit) (cid.Cid, error) {
	node, err := cbornode.WrapObject(commit, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := ds.Add(ctx, node); err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}

// CommitFromJSON decodes a commit from json.
func CommitFromJSON(data []byte) (*Commit, error) {
	var commit Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// CommitFromCBOR decodes a commit from an ipld node.
func CommitFromCBOR(data []byte) (*Commit, error) {
	var commit Commit
	if err := cbornode.DecodeInto(data, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// NewCommit returns a new commit with default values.
func NewCommit() *Commit {
	return &Commit{
		Date:     time.Now(),
		Metadata: make(map[string]string),
	}
}

// ParentLinks returns parent ipld links.
func (c *Commit) ParentLinks() []*ipld.Link {
	var out []*ipld.Link
	for _, p := range c.Parents {
		out = append(out, &ipld.Link{Cid: p})
	}

	return out
}
