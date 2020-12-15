package core

import (
	"context"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

// Worktree adds the current working tree to the merkle dag.
// Optional ignore rules can be used to filter out files.
func Worktree(ctx context.Context, fs afero.Fs, dag ipld.DAGService) (ipld.Node, error) {
	rules, err := Ignore(fs)
	if err != nil {
		return nil, err
	}

	filter, err := ignore.CompileIgnoreLines(rules...)
	if err != nil {
		return nil, err
	}

	return Add(ctx, fs, dag, "", filter)
}
