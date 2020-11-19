package object

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/util"
)

// Commit contains info about changes to a repo.
type Commit struct {
	// Date is the timestamp of when the commit was created.
	Date time.Time `refmt:"date,string"`
	// Message is a description of the changes.
	Message string `json:"message"`
	// Parents is a list of the parent commit CIDs.
	Parents []cid.Cid `json:"parents"`
	// Tree is the root CID of the repo file tree.
	Tree cid.Cid `json:"tree"`
}

// LogDateFormat is the format used when logging the commit.
const LogDateFormat = "Mon Jan 2 15:04:05 2006 -0700"

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
	out := []*ipld.Link{}
	for _, p := range c.Parents {
		out = append(out, &ipld.Link{Cid: p})
	}

	return out
}

// Log prints a human readable version of the commit.
func (c *Commit) Log(w io.Writer, id, head, base cid.Cid) {
	fmt.Fprintf(w, "%scommit %s", util.ColorYellow, id.String())

	if id == head {
		fmt.Fprintf(w, " (%sHEAD%s)", util.ColorRed, util.ColorYellow)
	}

	if id == base {
		fmt.Fprintf(w, " (%sBASE%s)", util.ColorGreen, util.ColorYellow)
	}

	fmt.Fprintf(w, "%s\n", util.ColorReset)
	fmt.Fprintf(w, "Date: %s\n", c.Date.Format(LogDateFormat))
	fmt.Fprintf(w, "\n\t%s\n\n", c.Message)
}
