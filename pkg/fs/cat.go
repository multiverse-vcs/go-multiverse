package fs

import (
	"context"
	"io/ioutil"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	io "github.com/ipfs/go-unixfs/io"
)

// Cat returns the contents of the file with the given CID.
func Cat(ctx context.Context, dag ipld.DAGService, id cid.Cid) (string, error) {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return "", err
	}

	reader, err := io.NewDagReader(ctx, node, dag)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
