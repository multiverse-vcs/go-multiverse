package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multihash"
)

var (
	// ErrInvalidCodec is returned when a cid codec is invalid.
	ErrInvalidCodec = errors.New("invalid codec type")
	// ErrNoLink is returned when an invalid path is resolved.
	ErrNoLink = errors.New("no such link")
	// ErrZeroLenPath is returned when a path of len zero is resolved.
	ErrZeroLenPath = errors.New("zero length path")
)

// Commit contains info about changes to a repo.
type Commit struct {
	// Date is the timestamp of when the commit was created.
	Date time.Time `json:"date"`
	// Message is a description of the changes.
	Message string `json:"message"`
	// Parents contains the CIDs of parent commits.
	Parents []cid.Cid `json:"parents"`
	// PeerID is the hash of the author's public key.
	PeerID peer.ID `json:"peer_id,omitempty"`
	// Worktree is the current state of the repo files.
	Worktree cid.Cid `json:"tree"`

	cid  cid.Cid
	data []byte
}

// Static (compile time) check that Commit satisfies the format.Node interface.
var _ format.Node = (*Commit)(nil)

// DecodeCommit decodes a commit from a cid and data.
func DecodeCommit(cid cid.Cid, data []byte) (*Commit, error) {
	if cid.Prefix().Codec != MCommit {
		return nil, ErrInvalidCodec
	}

	commit := Commit{cid: cid, data: data}
	if err := json.Unmarshal(data, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (c *Commit) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		return c, nil, nil
	}

	switch path[0] {
	case "date":
		return c.Date, path[1:], nil
	case "message":
		return c.Message, path[1:], nil
	case "peer_id":
		return c.PeerID, path[1:], nil
	case "tree":
		return &format.Link{Cid: c.Worktree}, path[1:], nil
	case "parents":
		if len(path) == 1 {
			return c.Parents, nil, nil
		}

		i, err := strconv.Atoi(path[1])
		if err != nil || i >= len(c.Parents) || i < 0 {
			return nil, nil, ErrNoLink
		}

		return &format.Link{Cid: c.Parents[i]}, path[2:], nil
	default:
		return nil, nil, ErrNoLink
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (c *Commit) Tree(path string, depth int) []string {
	if path != "" || depth == 0 {
		return nil
	}

	tree := []string{"date", "message", "peer_id", "tree"}
	for i := range c.Parents {
		tree = append(tree, fmt.Sprintf("parents/%d", i))
	}

	return tree
}

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
func (c *Commit) ResolveLink(path []string) (*format.Link, []string, error) {
	out, rest, err := c.Resolve(path)
	if err != nil {
		return nil, nil, err
	}

	lnk, ok := out.(*format.Link)
	if !ok {
		return nil, nil, ErrNoLink
	}

	return lnk, rest, nil
}

// Links is a helper function that returns all links within this object
func (c *Commit) Links() []*format.Link {
	out := []*format.Link{
		{Cid: c.Worktree},
	}

	for _, p := range c.Parents {
		out = append(out, &format.Link{Cid: p})
	}

	return out
}

// ParentLinks is a helper function that returns parent links.
func (c *Commit) ParentLinks() []*format.Link {
	out := []*format.Link{}

	for _, p := range c.Parents {
		out = append(out, &format.Link{Cid: p})
	}

	return out
}

// RawData returns the block raw contents as a byte slice.
func (c *Commit) RawData() []byte {
	if c.data != nil {
		return c.data
	}

	data, err := json.Marshal(c)
	if err != nil {
		panic("failed to marshal commit")
	}

	c.data = data
	return c.data
}

// Cid returns the cid of the commit.
func (c *Commit) Cid() cid.Cid {
	if c.cid.Defined() {
		return c.cid
	}

	hash, err := multihash.Sum(c.RawData(), multihash.SHA2_256, -1)
	if err != nil {
		panic("failed to hash commit")
	}

	c.cid = cid.NewCidV1(MCommit, hash)
	return c.cid
}

// Copy will go away. It is here to comply with the Node interface.
func (c *Commit) Copy() format.Node {
	return nil
}

// Size will go away. It is here to comply with the Node interface.
func (c *Commit) Size() (uint64, error) {
	return 0, nil
}

// Stat returns info about the node.
func (c *Commit) Stat() (*format.NodeStat, error) {
	return &format.NodeStat{}, nil
}

// String returns a string representation of the commit.
func (c *Commit) String() string {
	return "[multiverse commit]"
}

// Loggable returns a go-log loggable item.
func (c *Commit) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "multiverse-commit",
	}
}
