package p2p

import (
	"context"
	"errors"
	"path"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	namesys "github.com/libp2p/go-libp2p-pubsub-router"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Namespace is the pubsub topic namespace.
const Namespace = "multiverse"

// Pubsub allows subcribing to topics.
type Pubsub struct {
	store *namesys.PubsubValueStore
}

// TopicForPeerID returns the topic name for the given peer id.
func TopicForPeerID(id peer.ID) string {
	return path.Join("/", Namespace, string(id))
}

// NewPubsub returns a new pubsub router.
func NewPubsub(ctx context.Context, host host.Host) (*Pubsub, error) {
	sub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	store, err := namesys.NewPubsubValueStore(ctx, host, sub, validator{})
	if err != nil {
		return nil, err
	}

	return &Pubsub{
		store: store,
	}, nil
}

func (p *Pubsub) SearchAuthor(ctx context.Context, id peer.ID) (*data.Author, error) {
	out, err := p.store.SearchValue(ctx, TopicForPeerID(id))
	if err != nil {
		return nil, err
	}

	return data.AuthorFromCBOR(<-out)
}

func (p *Pubsub) GetAuthor(ctx context.Context, id peer.ID) (*data.Author, error) {
	value, err := p.store.GetValue(ctx, TopicForPeerID(id))
	if err != nil {
		return nil, err
	}

	return data.AuthorFromCBOR(value)
}

func (p *Pubsub) PutAuthor(ctx context.Context, id peer.ID, author *data.Author) error {
	node, err := author.Node()
	if err != nil {
		return err
	}

	return p.store.PutValue(ctx, TopicForPeerID(id), node.RawData())
}

type validator struct{}

func (v validator) Validate(key string, value []byte) error {
	ns, _, err := record.SplitKey(key)
	if err != nil {
		return err
	}

	if ns != Namespace {
		return errors.New("invalid namespace")
	}

	// TODO unmarshal author from value
	// TODO validate signature of author
	return nil
}

func (v validator) Select(key string, vals [][]byte) (int, error) {
	// TODO compare author sequence numbers
	return 0, nil
}
