package object

import (
	"context"
	"encoding/json"

	cid "github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multihash"
)

// Author contains info about a user.
type Author struct {
	// Repositories is a map of repositories.
	Repositories map[string]cid.Cid `json:"repositories"`
	// Following is a list of peer IDs to follow.
	Following []peer.ID `json:"following"`
	// Metadata contains additional data.
	Metadata map[string]string `json:"metadata"`
}

// GetAuthor returns the author with the given CID.
func GetAuthor(ctx context.Context, ds ipld.NodeGetter, id cid.Cid) (*Author, error) {
	node, err := ds.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return AuthorFromCBOR(node.RawData())
}

// AddAuthor adds a author to the given dag.
func AddAuthor(ctx context.Context, ds ipld.NodeAdder, author *Author) (cid.Cid, error) {
	node, err := cbornode.WrapObject(author, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := ds.Add(ctx, node); err != nil {
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
