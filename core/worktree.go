package core

import (
	"context"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/sabhiram/go-gitignore"
)

// Worktree adds the current working tree to the merkle dag.
// Optional ignore rules can be used to filter out files.
func Worktree(ctx context.Context, dag ipld.DAGService, path string) (ipld.Node, error) {
	rules, err := Ignore()
	if err != nil {
		return nil, err
	}

	filter, err := ignore.CompileIgnoreLines(rules...)
	if err != nil {
		return nil, err
	}

	return Add(ctx, dag, path, filter)
}
