package object

import (
	"context"
	"encoding/json"

	cid "github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

// Repository contains all versions of a project.
type Repository struct {
	// DefaultBranch is the base branch of the repo.
	DefaultBranch string `json:"default_branch"`
	// Branches is a map of names to commit CIDs.
	Branches map[string]cid.Cid `json:"branches"`
	// Tags is a map of names to commit CIDs.
	Tags map[string]cid.Cid `json:"tags"`
	// Metadata contains additional data.
	Metadata map[string]string `json:"metadata"`
}

// GetRepository returns the repo with the given CID.
func GetRepository(ctx context.Context, ds ipld.NodeGetter, id cid.Cid) (*Repository, error) {
	node, err := ds.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return RepositoryFromCBOR(node.RawData())
}

// AddRepository adds a repo to the given dag.
func AddRepository(ctx context.Context, ds ipld.NodeAdder, repo *Repository) (cid.Cid, error) {
	node, err := cbornode.WrapObject(repo, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := ds.Add(ctx, node); err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}

// RepositoryFromJSON decodes a repo from json.
func RepositoryFromJSON(data []byte) (*Repository, error) {
	var repo Repository
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

// RepositoryFromCBOR decodes a repo from an ipld node.
func RepositoryFromCBOR(data []byte) (*Repository, error) {
	var repo Repository
	if err := cbornode.DecodeInto(data, &repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

// NewRepository returns a new repo.
func NewRepository() *Repository {
	return &Repository{
		Branches: make(map[string]cid.Cid),
		Tags:     make(map[string]cid.Cid),
		Metadata: make(map[string]string),
	}
}

// Heads returns a set of all branch heads.
func (r *Repository) Heads() *cid.Set {
	heads := cid.NewSet()
	for _, id := range r.Branches {
		heads.Add(id)
	}
	return heads
}
