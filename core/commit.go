package core

import (
	"context"
	"time"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Commit creates a new commit.
func Commit(ctx context.Context, dag ipld.DAGService, path string, filter Filter, message string, parents ...cid.Cid) (cid.Cid, error) {
	tree, err := Add(ctx, dag, path, filter)
	if err != nil {
		return cid.Cid{}, err
	}

	commit := &data.Commit{
		Date:    time.Now(),
		Message: message,
		Tree:    tree.Cid(),
		Parents: parents,
	}

	node, err := commit.Node()
	if err != nil {
		return cid.Cid{}, err
	}

	if err := dag.Add(ctx, node); err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}
