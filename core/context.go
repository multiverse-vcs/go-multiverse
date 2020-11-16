package core

import (
	"context"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
)

// Context contains common data and services.
type Context struct {
	ctx context.Context
	dag ipld.DAGService
}

// NewMockContext returns a context that can be used for testing.
func NewMockContext() *Context {
	return &Context{
		ctx: context.TODO(),
		dag: dagutils.NewMemoryDagService(),
	}
}
