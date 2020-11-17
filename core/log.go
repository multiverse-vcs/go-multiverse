package core

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/util"
)

func (c *Context) Log(w io.Writer) error {
	if !c.config.Head.Defined() {
		return nil
	}

	getLinks := func(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
		node, err := c.dag.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		commit, ok := node.(*object.Commit)
		if !ok {
			return nil, errors.New("invalid commit")
		}

		return commit.ParentLinks(), nil
	}

	visit := func(id cid.Cid) bool {
		node, err := c.dag.Get(c.ctx, id)
		if err != nil {
			return false
		}

		commit, ok := node.(*object.Commit)
		if !ok {
			return false
		}

		fmt.Fprintf(w, "%scommit %s", util.ColorYellow, commit.Cid().String())
		if id == c.config.Head {
			fmt.Fprintf(w, " (%sHEAD%s)", util.ColorRed, util.ColorYellow)
		}
		if id == c.config.Base {
			fmt.Fprintf(w, " (%sBASE%s)", util.ColorGreen, util.ColorYellow)
		}
		fmt.Fprintf(w, "%s\n", util.ColorReset)
		fmt.Fprintf(w, "Peer: %s\n", commit.PeerID.String())
		fmt.Fprintf(w, "Date: %s\n", commit.Date.Format("Mon Jan 2 15:04:05 2006 -0700"))
		fmt.Fprintf(w, "\n\t%s\n\n", commit.Message)
		return true
	}

	return merkledag.Walk(c.ctx, getLinks, c.config.Head, visit)
}
