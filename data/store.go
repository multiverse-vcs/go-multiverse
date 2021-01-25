package data

import (
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"github.com/ipfs/go-datastore/query"
)

// prefix is the parent key for the datstore.
var prefix = datastore.NewKey("multiverse")

// Store is a key value database.
type Store struct {
	datastore.Datastore
}

// NewStore returns a new store that wraps the given store.
func NewStore(dstore datastore.Datastore) *Store {
	return &Store{namespace.Wrap(dstore, prefix)}
}

// GetCid returns the cid value of the given key name.
func (s *Store) GetCid(name string) (cid.Cid, error) {
	data, err := s.Get(datastore.NewKey(name))
	if err != nil {
		return cid.Cid{}, err
	}

	return cid.Cast(data)
}

// PutCid persists the given key name and value id.
func (s *Store) PutCid(name string, id cid.Cid) error {
	return s.Put(datastore.NewKey(name), id.Bytes())
}

// Keys returns a list of all keys in the store.
func (s *Store) Keys() ([]string, error) {
	res, err := s.Query(query.Query{KeysOnly: true})
	if err != nil {
		return nil, err
	}

	all, err := res.Rest()
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, e := range all {
		keys = append(keys, e.Key)
	}

	return keys, nil
}