package unixfs

import (
	"context"
	"strings"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/nasdf/diff3"
)

// Merge combines the contents of two edited files into the original.
func Merge(ctx context.Context, dag ipld.DAGService, o, a, b cid.Cid) (ipld.Node, error) {
	textO, err := Cat(ctx, dag, o)
	if err != nil {
		return nil, err
	}

	textA, err := Cat(ctx, dag, a)
	if err != nil {
		return nil, err
	}

	textB, err := Cat(ctx, dag, b)
	if err != nil {
		return nil, err
	}

	merged := diff3.Merge(textO, textA, textB)
	reader := strings.NewReader(merged)

	return addReader(ctx, dag, reader)
}
