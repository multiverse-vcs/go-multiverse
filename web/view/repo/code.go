package repo

import (
	"context"
	"regexp"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-path"
	"github.com/ipfs/go-unixfs"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
)

var readmeRegex = regexp.MustCompile(`(?i)^readme.*`)

// CodeModel contains data for the code subview.
type CodeModel struct {
	Blob   string
	IsDir  bool
	Readme string
	Tree   []*core.DirEntry
}

// Code returns a new code model.
func Code(ctx context.Context, node *node.Node, id cid.Cid, file string) (*CodeModel, error) {
	fpath, err := path.FromSegments("/ipfs/", id.String(), "tree", file)
	if err != nil {
		return nil, err
	}

	fnode, err := node.ResolvePath(ctx, fpath)
	if err != nil {
		return nil, err
	}

	fsnode, err := unixfs.ExtractFSNode(fnode)
	if err != nil {
		return nil, err
	}

	model := CodeModel{
		IsDir: fsnode.IsDir(),
	}

	switch {
	case model.IsDir:
		tree, err := core.Ls(ctx, node, fnode.Cid())
		if err != nil {
			return nil, err
		}

		readme, err := readme(ctx, node, tree)
		if err != nil {
			return nil, err
		}

		model.Readme = readme
		model.Tree = tree
	default:
		blob, err := core.Cat(ctx, node, fnode.Cid())
		if err != nil {
			return nil, err
		}

		model.Blob = blob
	}

	return &model, nil
}

// readme returns the contents of the readme if it exists.
func readme(ctx context.Context, node *node.Node, tree []*core.DirEntry) (string, error) {
	for _, e := range tree {
		if readmeRegex.MatchString(e.Name) {
			return core.Cat(ctx, node, e.Cid)
		}
	}
	return "", nil
}
