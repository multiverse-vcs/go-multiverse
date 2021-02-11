package data

import (
	"context"
	"encoding/json"

	"github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multihash"
)

// Author contains info about a user.
type Author struct {
	// Repositories is a map of repositories.
	Repositories map[string]cid.Cid `json:"repositories"`
	// Metadata contains additional data.
	Metadata map[string]string `json:"metadata"`
	// Following contains a list of followed peers.
	Following []peer.ID `json:"following"`
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
		Metadata:     make(map[string]string),
	}
}
