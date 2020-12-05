package core

import (
	"context"
	"io/ioutil"

	"github.com/ipfs/go-cid"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// Cat returns the contents of a file.
func Cat(ctx context.Context, store *storage.Store, id cid.Cid) (string, error) {
	node, err := store.Dag.Get(ctx, id)
	if err != nil {
		return "", err
	}

	reader, err := ufsio.NewDagReader(ctx, node, store.Dag)
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
