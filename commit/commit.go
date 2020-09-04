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

// Add creates a commit from local changes.
func Add(ipfs *core.IpfsNode, message string, changes cid.Cid, parent cid.Cid) (format.Node, error) {
	node := &Commit{Message: message, Changes: changes, Parent: parent}

	dag, err := cbornode.WrapObject(node, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	if err := ipfs.DAG.Add(context.TODO(), dag); err != nil {
		return nil, err
	}

	return dag, nil
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

// Log prints the commit history starting at the given CID.
func Log(ipfs *core.IpfsNode, id cid.Cid) error {
	c, err := Get(ipfs, id)
	if err != nil {
		return err
	}

	fmt.Println(c.String())
	if c.Parent.Defined() {
		return Log(ipfs, c.Parent)
	}

	return nil
}

// String returns a human readable representation of the commit.
func (c *Commit) String() string {
	return fmt.Sprintf("%s\n\t%s\n", c.Changes, c.Message)
}
