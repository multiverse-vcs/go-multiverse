package data

import (
	"context"
	"encoding/json"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-pinner"
	"github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multihash"
)

// Repository contains all versions of a project.
type Repository struct {
	// Author is the peer id of the author.
	Author peer.ID `json:"author"`
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

// GetRepository returns the repo with the given CID.
func GetRepository(ctx context.Context, dag ipld.DAGService, id cid.Cid) (*Repository, error) {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return RepositoryFromCBOR(node.RawData())
}

// AddRepository adds a repo to the given dag.
func AddRepository(ctx context.Context, dag ipld.DAGService, repo *Repository) (cid.Cid, error) {
	node, err := cbornode.WrapObject(repo, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := dag.Add(ctx, node); err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}

// PinRepository pins a repo using the given pinner.
func PinRepository(ctx context.Context, pinner pin.Pinner, repo *Repository) (cid.Cid, error) {
	node, err := cbornode.WrapObject(repo, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := pinner.Pin(ctx, node, true); err != nil {
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
func NewRepository(name string) *Repository {
	return &Repository{
		Name:     name,
		Branches: make(map[string]cid.Cid),
		Tags:     make(map[string]cid.Cid),
		Metadata: make(map[string]string),
	}
}

// Ref returns the cid of the given ref.
func (r *Repository) Ref(ref string) (cid.Cid, error) {
	if id, ok := r.Branches[ref]; ok {
		return id, nil
	}

	if id, ok := r.Tags[ref]; ok {
		return id, nil
	}

	return cid.Parse(ref)
}
