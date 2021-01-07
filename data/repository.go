package data

import (
	"encoding/json"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

// Repository contains all versions of a project.
type Repository struct {
	// Name is the human friendly name of the repo.
	Name string `json:"name"`
	// Description describes the project.
	Description string `json:"description"`
	// Branches is a map of names to commit CIDs.
	Branches map[string]cid.Cid `json:"branches"`
	// Tags is a map of names to commit CIDs.
	Tags map[string]cid.Cid `json:"tags"`
	// Metadata contains additional data.
	Metadata map[string]string `json:"metadata"`
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
func NewRepository(name string) *Repository {
	return &Repository{
		Name:     name,
		Branches: make(map[string]cid.Cid),
		Tags:     make(map[string]cid.Cid),
		Metadata: make(map[string]string),
	}
}

// Node returns an ipld node containing the commit.
func (r *Repository) Node() (ipld.Node, error) {
	return cbornode.WrapObject(&r, multihash.SHA2_256, -1)
}

// DefaultBranch returns the default branch for the repo.
func (r *Repository) DefaultBranch() string {
	if _, ok := r.Branches["default"]; ok {
		return "default"
	}

	for k := range r.Branches {
		return k
	}

	return ""
}
