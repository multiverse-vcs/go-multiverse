package core

import (
	"context"
	"io/ioutil"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
)

// Config contains common configuration info.
type Config struct {
	Base cid.Cid
	Head cid.Cid
}

// Context contains common data and services.
type Context struct {
	ctx    context.Context
	config *Config
	dag    ipld.DAGService
	root   string
}

// NewMockContext returns a context that can be used for testing.
func NewMockContext() (*Context, error) {
	root, err := ioutil.TempDir("", "multiverse-*")
	if err != nil {
		return nil, err
	}

	c := Context{
		ctx:    context.TODO(),
		config: &Config{},
		dag:    dagutils.NewMemoryDagService(),
		root:   root,
	}

	return &c, nil
}
