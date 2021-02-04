package peer

import (
	"github.com/ipfs/go-datastore"
)

const megabyte = 1024 * 1024

// Metrics contains info about resource usage.
type Metrics struct {
	// DiskUsage is the amount of storage used.
	DiskUsage uint64
	// Peers is the number of discovered peers.
	Peers int
}

// GetMetrics returns a snapshot of the current metrics.
func (c *Client) GetMetrics() (*Metrics, error) {
	peers := c.host.Peerstore().Peers()

	diskUsage, err := datastore.DiskUsage(c.dstore)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		DiskUsage: diskUsage / megabyte,
		Peers:     len(peers),
	}, nil
}
