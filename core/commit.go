package core

import (
	"context"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// Commit creates a new commit.
func Commit(ctx context.Context, store *storage.Store, message string, parents ...cid.Cid) (cid.Cid, error) {
	tree, err := Worktree(ctx, store)
	if err != nil {
		return cid.Cid{}, err
	}

	commit := object.Commit{
		Date:    time.Now(),
		Message: message,
		Tree:    tree.Cid(),
		Parents: parents,
	}

	node, err := cbornode.WrapObject(&commit, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := store.Dag.Add(ctx, node); err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}
