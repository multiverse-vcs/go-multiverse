package name

import (
	"context"
	"path"

	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	namesys "github.com/libp2p/go-libp2p-pubsub-router"
)

// Namespace is the pubsub topic namespace.
const Namespace = "multiverse"

// System performs name resolution.
type System struct {
	values *namesys.PubsubValueStore
}

// TopicForPeerID returns the topic name for the given peer id.
func TopicForPeerID(id peer.ID) string {
	return path.Join("/", Namespace, peer.Encode(id))
}

// NewNameSystem returns a new name system.
func NewSystem(ctx context.Context, host host.Host, router routing.Routing, dstore datastore.Datastore) (*System, error) {
	dis := discovery.NewRoutingDiscovery(router)

	// TODO use datastore to persist values

	sub, err := pubsub.NewGossipSub(ctx, host, pubsub.WithDiscovery(dis))
	if err != nil {
		return nil, err
	}

	values, err := namesys.NewPubsubValueStore(ctx, host, sub, Validator{})
	if err != nil {
		return nil, err
	}

	return &System{
		values: values,
	}, nil
}

// GetValue returns the latest value for the topic with the given peer id.
func (s *System) GetValue(ctx context.Context, id peer.ID) (*Record, error) {
	val, err := s.values.GetValue(ctx, TopicForPeerID(id))
	if err != nil {
		return nil, err
	}

	return RecordFromCBOR(val)
}

// PutValue publishes the value under the topic of the given peer id.
func (s *System) PutValue(ctx context.Context, id peer.ID, rec *Record) error {
	val, err := rec.Bytes()
	if err != nil {
		return err
	}

	return s.values.PutValue(ctx, TopicForPeerID(id), val)
}

// Search searches for the the latest value from the topic with the given peer ID.
func (s *System) SearchValue(ctx context.Context, id peer.ID) (*Record, error) {
	out, err := s.values.SearchValue(ctx, TopicForPeerID(id))
	if err != nil {
		return nil, err
	}

	val, ok := <-out
	if !ok {
		return nil, routing.ErrNotFound
	}

	return RecordFromCBOR(val)
}

// Subscribe creates a subscription to the topic of the given peer ID.
func (s *System) Subscribe(id peer.ID) error {
	return s.values.Subscribe(TopicForPeerID(id))
}

// Unsubscribe cancels a subscription to the topic of the given peer ID.
func (s *System) Unsubscribe(id peer.ID) (bool, error) {
	return s.values.Cancel(TopicForPeerID(id))
}

// Publish advertises the given id to the topic of the peer ID from the private key.
func (s *System) Publish(ctx context.Context, key crypto.PrivKey, id cid.Cid) error {
	peerID, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return err
	}

	val, err := s.GetValue(ctx, peerID)
	if err != nil && err != routing.ErrNotFound {
		return err
	}

	rec := NewRecord(id.Bytes())
	if val != nil {
		rec.Sequence = val.Sequence + 1
	}

	if err := rec.Sign(key); err != nil {
		return err
	}

	return s.PutValue(ctx, peerID, rec)
}

// Resolve returns the latest value from the topic with the given peer ID.
func (s *System) Resolve(ctx context.Context, id peer.ID) (cid.Cid, error) {
	rec, err := s.GetValue(ctx, id)
	if err != nil {
		return cid.Cid{}, err
	}

	return cid.Cast(rec.Value)
}

// Search searches for the the latest value from the topic with the given peer ID.
func (s *System) Search(ctx context.Context, id peer.ID) (cid.Cid, error) {
	rec, err := s.SearchValue(ctx, id)
	if err != nil {
		return cid.Cid{}, err
	}

	return cid.Cast(rec.Value)
}