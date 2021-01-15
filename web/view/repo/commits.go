package repo

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/node"
)

// CommitsModel contains data for the commits subview.
type CommitsModel struct {
	IDs  []cid.Cid
	List []*data.Commit
}

// Commits returns a new commits model.
func Commits(ctx context.Context, node *node.Node, id cid.Cid) (*CommitsModel, error) {
	var ids []cid.Cid
	visit := func(id cid.Cid) bool {
		ids = append(ids, id)
		return true
	}

	_, err := core.Walk(ctx, node, id, visit)
	if err != nil {
		return nil, err
	}

	var list []*data.Commit
	for _, id := range ids {
		commit, err := data.GetCommit(ctx, node, id)
		if err != nil {
			return nil, err
		}

		list = append(list, commit)
	}

	return &CommitsModel{
		IDs:  ids,
		List: list,
	}, nil
}
