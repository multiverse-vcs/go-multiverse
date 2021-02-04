package unixfs

import (
	"context"
	"regexp"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	ufs "github.com/ipfs/go-unixfs"
)

// IsDir returns true if the file with the given id is a unixfs directory.
func IsDir(ctx context.Context, dag ipld.DAGService, id cid.Cid) (bool, error) {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return false, err
	}

	fsnode, err := ufs.ExtractFSNode(node)
	if err != nil {
		return false, err
	}

	return fsnode.IsDir(), nil
}

// Find returns the first dir entry matching the given pattern in the directory with the given id.
func Find(ctx context.Context, dag ipld.DAGService, id cid.Cid, pattern *regexp.Regexp) (*DirEntry, error) {
	entries, err := Ls(ctx, dag, id)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if pattern.MatchString(e.Name) {
			return e, nil
		}
	}

	return nil, nil
}
