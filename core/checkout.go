package core

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/spf13/afero"
)

// Checkout writes the tree of the commit to the root.
func Checkout(ctx context.Context, fs afero.Fs, dag ipld.DAGService, id cid.Cid) error {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return errors.New("invalid commit")
	}

	tree, err := dag.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	return Write(ctx, fs, dag, "", tree)
}
