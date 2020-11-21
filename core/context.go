package core

import (
	"context"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/config"
)

// Context contains common data and services.
type Context struct {
	context.Context
	cfg *config.Config
	dag ipld.DAGService
	fs  billy.Filesystem
}

// NewMockContext returns a context that can be used for testing.
func NewMockContext() *Context {
	return &Context{
		Context: context.TODO(),
		cfg:     &config.Config{},
		dag:     dagutils.NewMemoryDagService(),
		fs:      memfs.New(),
	}
}
