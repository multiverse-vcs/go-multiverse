package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
)

// Parents contains parent CIDs.
type Parents []cid.Cid

// Signature contains info about who created the commit.
type Signature struct {
	// Name is the name of the person who created the commit.
	Name string `json:"name"`
	// Email is an address that can be used to contact the committer.
	Email string `json:"email"`
	// Date is the timestamp of when the commit was created.
	Date time.Time `json:"date"`
}

// Commit contains info about changes to a repo.
type Commit struct {
	// Author is the person that created the commit.
	Author *Signature `json:"author"`
	// Committer is the person that performed the commit.
	Committer *Signature `json:"committer"`
	// Message is a description of the changes.
	Message string `json:"message"`
	// WorkTree is the current state of the repo files.
	WorkTree cid.Cid `json:"tree"`
	// Parents contains the CIDs of parent commits.
	Parents Parents `json:"parents"`

	cid  cid.Cid
	data []byte
}

// Static (compile time) check that Commit satisfies the format.Node interface.
var _ format.Node = (*Commit)(nil)

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (s *Signature) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		return s, nil, nil
	}

	switch path[0] {
	case "name":
		return s.Name, path[1:], nil
	case "email":
		return s.Email, path[1:], nil
	case "date":
		return s.Date, path[1:], nil
	default:
		return nil, nil, ErrNoLink
	}
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (p Parents) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		return p, nil, nil
	}

	i, err := strconv.Atoi(path[0])
	if err != nil {
		return nil, nil, err
	}

	if i < 0 || i >= len(p) {
		return nil, nil, ErrOutOfRange
	}

	return &format.Link{Cid: p[i]}, path[1:], nil
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (c *Commit) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		return c, nil, nil
	}

	switch path[0] {
	case "message":
		return c.Message, path[1:], nil
	case "author":
		return c.Author.Resolve(path[1:])
	case "committer":
		return c.Committer.Resolve(path[1:])
	case "parents":
		return c.Parents.Resolve(path[1:])
	case "tree":
		return &format.Link{Cid: c.WorkTree}, path[1:], nil
	}

	return nil, nil, ErrNoLink
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (s *Signature) Tree(path string, depth int) []string {
	if path != "" || depth == 0 {
		return nil
	}

	tree := []string{"name", "email", "date"}
	for i := range tree {
		tree[i] = fmt.Sprintf("%s/%s", path, tree[i])
	}

	return tree
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (p Parents) Tree(path string, depth int) []string {
	if path != "" || depth == 0 {
		return nil
	}

	tree := make([]string, len(p))
	for i := range p {
		tree[i] = fmt.Sprintf("%s/%d", path, i)
	}

	return tree
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (c *Commit) Tree(path string, depth int) []string {
	if path != "" || depth == 0 {
		return nil
	}

	tree := []string{"tree", "message"}
	tree = append(tree, c.Author.Tree("author", depth)...)
	tree = append(tree, c.Committer.Tree("committer", depth)...)
	tree = append(tree, c.Parents.Tree("parents", depth)...)
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

// Copy returns a deep copy of this node
func (c* Commit) Copy() format.Node {
	nc := *c
	return &nc
}

// Links is a helper function that returns all links within this object
func (c *Commit) Links() []*format.Link {
	out := make([]*format.Link, len(c.Parents) + 1)

	for _, p := range c.Parents {
		out = append(out, &format.Link{Cid: p})
	}

	return append(out, &format.Link{Cid: c.WorkTree})
}

// Stat returns info about the node.
func (c *Commit) Stat() (*format.NodeStat, error) {
	return &format.NodeStat{}, nil
}

// Size returns the size in bytes of the serialized object
func (c *Commit) Size() (uint64, error) {
	size := len(c.RawData())
	return uint64(size), nil
}

// RawData returns the block raw contents as a byte slice.
func (c *Commit) RawData() []byte {
	return c.data
}

// Cid returns the cid of the commit.
func (c *Commit) Cid() cid.Cid {
	return c.cid
}

// String returns a string representation of the commit.
func (c *Commit) String() string {
	return "[multiverse commit]"
}

// Loggable returns a go-log loggable item.
func (c *Commit) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "multi-commit",
	}
}