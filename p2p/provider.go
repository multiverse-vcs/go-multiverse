package p2p

import (
	"context"
	"time"

	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	provider "github.com/ipfs/go-ipfs-provider"
	"github.com/ipfs/go-ipfs-provider/queue"
	"github.com/ipfs/go-ipfs-provider/simple"
	"github.com/libp2p/go-libp2p-core/routing"
)

const (
	// ReprovideInterval is the time between reprovides.
	ReprovideInterval = 12 * time.Hour
	// QueueName is the name for the provider queue.
	QueueName = "repro"
)

// NewProvider returns a new provider system.
func NewProvider(ctx context.Context, ds datastore.Datastore, bs blockstore.Blockstore, router routing.Routing) (provider.System, error) {
	queue, err := queue.NewQueue(ctx, QueueName, ds)
	if err != nil {
		return nil, err
	}

	prov := simple.NewProvider(ctx, queue, router)
	bspr := simple.NewBlockstoreProvider(bs)
	repr := simple.NewReprovider(ctx, ReprovideInterval, router, bspr)

	return provider.NewSystem(prov, repr), nil
}
