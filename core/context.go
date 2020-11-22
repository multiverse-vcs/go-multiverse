package core

import (
	"context"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/config"
)

// Context is the context for core commands.
type Context struct {
	context.Context
	Config *config.Config
	Dag    ipld.DAGService
	Fs     billy.Filesystem
}

// NewMockContext returns a context that can be used for testing.
func NewMockContext() *Context {
	return &Context{
		Context: context.TODO(),
		Config:  &config.Config{},
		Dag:     dagutils.NewMemoryDagService(),
		Fs:      memfs.New(),
	}
}
