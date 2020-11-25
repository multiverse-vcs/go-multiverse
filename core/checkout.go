package core

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// Checkout writes the tree of the commit to the root.
func Checkout(ctx context.Context, store *storage.Store, id cid.Cid) error {
	node, err := store.Dag.Get(ctx, id)
	if err != nil {
		return err
	}

	commit, err := object.CommitFromCBOR(node.RawData())
	if err != nil {
		return errors.New("invalid commit")
	}

	tree, err := store.Dag.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	return Write(ctx, store, "", tree)
}
