package dag

import (
	"context"
	"io"

	chunk "github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs/importer/balanced"
	"github.com/ipfs/go-unixfs/importer/helpers"
)

// DefaultChunker is the name of the default chunker algorithm.
const DefaultChunker = "buzhash"

// Chunk splits the given reader into chunks and arranges them into a dag node.
func Chunk(ctx context.Context, dag ipld.DAGService, reader io.Reader) (ipld.Node, error) {
	chunker, err := chunk.FromString(reader, DefaultChunker)
	if err != nil {
		return nil, err
	}

	params := helpers.DagBuilderParams{
		Dagserv:    dag,
		CidBuilder: merkledag.V1CidPrefix(),
		Maxlinks:   helpers.DefaultLinksPerBlock,
	}

	helper, err := params.New(chunker)
	if err != nil {
		return nil, err
	}

	node, err := balanced.Layout(helper)
	if err != nil {
		return nil, err
	}

	if err := dag.Add(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}
