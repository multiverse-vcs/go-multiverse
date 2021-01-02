package node

import (
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

// repo is used to access data objects.
type repo struct {
	dstore datastore.Batching
}

// Set sets the head of repo with the given name.
func (r *repo) Set(name string, id cid.Cid) error {
	return r.dstore.Put(datastore.NewKey(name), id.Bytes())
}

// Get returns the head of the repo with the given name.
func (r *repo) Get(name string) (cid.Cid, error) {
	val, err := r.dstore.Get(datastore.NewKey(name))
	if err != nil {
		return cid.Cid{}, err
	}

	return cid.Cast(val)
}

// List returns a list of all repos.
func (r *repo) List() ([]query.Entry, error) {
	res, err := r.dstore.Query(query.Query{})
	if err != nil {
		return nil, err
	}

	return res.Rest()
}
