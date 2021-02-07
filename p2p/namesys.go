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

// TopicForPeerID returns the topic name for the given peer id.
func TopicForPeerID(id peer.ID) string {
	return path.Join("/", Namespace, peer.Encode(id))
}

// PeerIDForTopic returns the peer id for the given topic.
func PeerIDForTopic(topic string) (peer.ID, error) {
	ns, k, err := record.SplitKey(topic)
	if err != nil {
		return peer.ID(""), err
	}

	if ns != Namespace {
		return peer.ID(""), record.ErrInvalidRecordType
	}

	return peer.Decode(k)
}

// NewSystem returns a new pubsub name system.
func NewNamesys(ctx context.Context, host host.Host) (*namesys.PubsubValueStore, error) {
	sub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	return namesys.NewPubsubValueStore(ctx, host, sub, validator{})
}

type validator struct{}

// Validate ensures that the signature matches the topic id.
func (v validator) Validate(key string, value []byte) error {
	id, err := PeerIDForTopic(key)
	if err != nil {
		return err
	}

	pub, err := id.ExtractPublicKey()
	if err != nil {
		return err
	}

	rec, err := data.RecordFromCBOR(value)
	if err != nil {
		return err
	}

	match, err := rec.Verify(pub)
	if err != nil {
		return err
	}

	if !match {
		return errors.New("signature does not match")
	}

	return nil
}

// Select finds the best record by comparing sequence numbers.
func (v validator) Select(key string, vals [][]byte) (int, error) {
	ind, max := -1, uint64(0)

	for i, v := range vals {
		rec, err := data.RecordFromCBOR(v)
		if err != nil {
			return -1, err
		}

		if rec.Sequence > max {
			ind, max = i, rec.Sequence
		}
	}

	return ind, nil
}
