package core

import (
	"context"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/sabhiram/go-gitignore"
)

// Worktree adds the current working tree to the merkle dag.
// Optional ignore rules can be used to filter out files.
func Worktree(ctx context.Context, store *storage.Store) (ipld.Node, error) {
	rules, err := Ignore(store)
	if err != nil {
		return nil, err
	}

	filter, err := ignore.CompileIgnoreLines(rules...)
	if err != nil {
		return nil, err
	}

	return Add(ctx, store, "", filter)
}
