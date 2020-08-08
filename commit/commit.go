package commit

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/multiformats/go-multihash"
)

type Commit struct {
	Id      cid.Cid `refmt:"-"`
	Message string  `refmt:"message,omitempty"`
	Changes cid.Cid `refmt:"changes,omitempty"`
	Parent  cid.Cid `refmt:"parent,omitempty"`
}

func init() {
	cbornode.RegisterCborType(Commit{})
}

// Creates a new commit with the message, changes, and parent.
func NewCommit(message string, changes cid.Cid, parent cid.Cid) *Commit {
	return &Commit{Message: message, Changes: changes, Parent: parent}
}

// Get the commit with the matching CID from IPFS.
func Get(ipfs iface.CoreAPI, id cid.Cid) (*Commit, error) {
	dag, err := ipfs.Dag().Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	var com Commit
	if err := cbornode.DecodeInto(dag.RawData(), &com); err != nil {
		return nil, err
	}

	com.Id = id
	return &com, nil
}

// Return a cbor node representation of the commit.
func (c *Commit) Node(ipfs iface.CoreAPI) (*cbornode.Node, error) {
	return cbornode.WrapObject(c, multihash.SHA2_256, -1)
}

// Return a human readable representation of the commit.
func (c *Commit) String() string {
	return fmt.Sprintf("%s\n%s", c.Message, c.Changes)
}