package core

import (
	"context"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
)

// Config contains common configuration info.
type Config struct {
	Root string  `json:"-"`
	Base cid.Cid `json:"base"`
	Head cid.Cid `json:"head"`
}

// Context contains common data and services.
type Context struct {
	ctx    context.Context
	config *Config
	dag    ipld.DAGService
	fs     billy.Filesystem
}

// NewMockContext returns a context that can be used for testing.
func NewMockContext() *Context {
	return &Context{
		ctx:    context.TODO(),
		config: &Config{},
		dag:    dagutils.NewMemoryDagService(),
		fs:     memfs.New(),
	}
}
