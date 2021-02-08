package peer

import (
	"context"
	"errors"

	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// AuthorsAPI implements methods to manage authors.
type AuthorsAPI struct {
	Peer
}

// Publish advertises the local author.
func (a *AuthorsAPI) Publish(ctx context.Context) error {
	config := a.Config()

	key, err := p2p.DecodeKey(config.PrivateKey)
	if err != nil {
		return err
	}

	id, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return err
	}

	payload, err := cbornode.DumpObject(config.Author)
	if err != nil {
		return err
	}

	signature, err := key.Sign(payload)
	if err != nil {
		return err
	}

	rec := data.NewRecord(payload, config.Sequence, signature)

	val, err := cbornode.DumpObject(rec)
	if err != nil {
		return err
	}

	return a.Namesys().PutValue(ctx, p2p.TopicForPeerID(id), val)
}

// Search returns the author published under the given peer id.
func (a *AuthorsAPI) Search(ctx context.Context, id peer.ID) (*data.Author, error) {
	out, err := a.Namesys().SearchValue(ctx, p2p.TopicForPeerID(id))
	if err != nil {
		return nil, err
	}

	val, ok := <-out
	if !ok {
		return nil, errors.New("author not found")
	}

	rec, err := data.RecordFromCBOR(val)
	if err != nil {
		return nil, err
	}

	return data.AuthorFromCBOR(rec.Payload)
}
