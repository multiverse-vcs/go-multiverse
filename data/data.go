// Package data contains object definitions.
package data

import (
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"github.com/ipfs/go-datastore/query"
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/polydawn/refmt/obj/atlas"
)

// timeAtlasEntry allows encoding and decoding of time structs.
var timeAtlasEntry = atlas.BuildEntry(time.Time{}).
	Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(t time.Time) (string, error) {
			return t.Format(time.RFC3339), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(t string) (time.Time, error) {
			return time.Parse(time.RFC3339, t)
		})).
	Complete()

func init() {
	cbornode.RegisterCborType(timeAtlasEntry)
	cbornode.RegisterCborType(Commit{})
	cbornode.RegisterCborType(Repository{})
}

// StorePrefix is the datastore key prefix.
var StorePrefix = datastore.NewKey("multiverse")

// Store is used to manage multiverse data.
type Store struct {
	datastore.Datastore
}

// NewStore returns a store backed by the given datastore.
func NewStore(dstore datastore.Datastore) *Store {
	return &Store{
		Datastore: namespace.Wrap(dstore, StorePrefix),
	}
}

// PutCid stores the given id under the given name.
func (s *Store) PutCid(name string, id cid.Cid) error {
	return s.Put(datastore.NewKey(name), id.Bytes())
}

// GetCid returns the cid with the given name.
func (s *Store) GetCid(name string) (cid.Cid, error) {
	val, err := s.Get(datastore.NewKey(name))
	if err != nil {
		return cid.Cid{}, err
	}

	return cid.Cast(val)
}

// Keys returns a list of all keys.
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
