package data

import (
	"context"
	"encoding/json"

	"github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

// Author contains info about a user.
type Author struct {
	// Name is the human friendly name of the author.
	Name string `json:"name"`
	// Email is the email address of the author.
	Email string `json:"email"`
	// Repositories is a map of repositories.
	Repositories map[string]cid.Cid `json:"repositories"`
}

// GetAuthor returns the author with the given CID.
func GetAuthor(ctx context.Context, dag ipld.DAGService, id cid.Cid) (*Author, error) {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return AuthorFromCBOR(node.RawData())
}

// AddAuthor adds a author to the given dag.
func AddAuthor(ctx context.Context, dag ipld.DAGService, author *Author) (cid.Cid, error) {
	node, err := cbornode.WrapObject(author, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := dag.Add(ctx, node); err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}

// AuthorFromJSON decodes a author from json.
func AuthorFromJSON(data []byte) (*Author, error) {
	var author Author
	if err := json.Unmarshal(data, &author); err != nil {
		return nil, err
	}

	return &author, nil
}

// AuthorFromCBOR decodes a author from an ipld node.
func AuthorFromCBOR(data []byte) (*Author, error) {
	var author Author
	if err := cbornode.DecodeInto(data, &author); err != nil {
		return nil, err
	}

	return &author, nil
}

// NewAuthor returns a new author.
func NewAuthor() *Author {
	return &Author{
		Repositories: make(map[string]cid.Cid),
	}
}
