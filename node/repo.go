package node

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// PutRepository stores the given repo.
func (n *Node) PutRepository(ctx context.Context, repo *data.Repository) error {
	id, err := data.AddRepository(ctx, n, repo)
	if err != nil {
		return err
	}

	return n.dstore.Put(datastore.NewKey(repo.Name), id.Bytes())
}

// GetRepository returns the repo with the given name.
func (n *Node) GetRepository(ctx context.Context, name string) (*data.Repository, error) {
	val, err := n.dstore.Get(datastore.NewKey(name))
	if err != nil {
		return nil, err
	}

	id, err := cid.Cast(val)
	if err != nil {
		return nil, err
	}

	return data.GetRepository(ctx, n, id)
}

// GetRepositoryOrDefault returns the repo with the given name or one with default settings.
func (n *Node) GetRepositoryOrDefault(ctx context.Context, name string) (*data.Repository, error) {
	exists, err := n.dstore.Has(datastore.NewKey(name))
	if err != nil {
		return nil, err
	}

	if !exists {
		return data.NewRepository(name), nil
	}

	return n.GetRepository(ctx, name)
}

// ListRepositories returns a list of all repositories.
func (n *Node) ListRepositories(ctx context.Context) ([]*data.Repository, error) {
	res, err := n.dstore.Query(query.Query{KeysOnly: true})
	if err != nil {
		return nil, err
	}

	all, err := res.Rest()
	if err != nil {
		return nil, err
	}

	var list []*data.Repository
	for _, e := range all {
		repo, err := n.GetRepository(ctx, e.Key)
		if err != nil {
			return nil, err
		}

		list = append(list, repo)
	}

	return list, nil
}
