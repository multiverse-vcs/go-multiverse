package merge

import (
	"context"
	"strings"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/nasdf/diff3"
)

// File combines the contents of two edited files into the original.
func File(ctx context.Context, ds ipld.DAGService, o, a, b cid.Cid) (ipld.Node, error) {
	textO, err := fs.Cat(ctx, ds, o)
	if err != nil {
		return nil, err
	}

	textA, err := fs.Cat(ctx, ds, a)
	if err != nil {
		return nil, err
	}

	textB, err := fs.Cat(ctx, ds, b)
	if err != nil {
		return nil, err
	}

	merged := diff3.Merge(textO, textA, textB)
	reader := strings.NewReader(merged)

	return dag.Chunk(ctx, ds, reader)
}
