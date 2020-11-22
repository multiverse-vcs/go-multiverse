package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/core"
)

const (
	// RootDir is the name of the root directory.
	RootDir = ".multiverse"
	// DatastoreName is the name of the datastore directory.
	DatastoreDir = "datastore"
)

func init() {
	core.IgnoreRules = append(core.IgnoreRules, RootDir)
}

// Context is the context for cli commands.
type Context struct {
	core.Context
}

// NewContext is used internally to create a command context.
func NewContext(root string, ctx context.Context) (*Context, error) {
	path := filepath.Join(root, RootDir, DatastoreDir)

	dstore, err := badger.NewDatastore(path, &badger.DefaultOptions)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(dstore)
	bserv := blockservice.New(bstore, offline.Exchange(bstore))

	corectx := core.Context{
		Context: ctx,
		Config:  &config.Config{},
		Dag:     merkledag.NewDAGService(bserv),
		Fs:      osfs.New(root),
	}

	return &Context{corectx}, nil
}

// InitContext initializes a context in the given path.
func InitContext(path string, ctx context.Context) (*Context, error) {
	if _, err := LoadContext(path, ctx); err == nil {
		return nil, errors.New("repo already exists")
	}

	root := filepath.Join(path, RootDir)
	if err := os.Mkdir(root, 0755); err != nil {
		return nil, err
	}

	return NewContext(path, ctx)
}

// LoadContext loads a context in the given path or parent directories.
func LoadContext(path string, ctx context.Context) (*Context, error) {
	root := filepath.Join(path, RootDir)

	info, err := os.Lstat(root)
	if err == nil && info.IsDir() {
		return NewContext(path, ctx)
	}

	parent := filepath.Dir(path)
	if parent == path {
		return nil, errors.New("repo not found")
	}

	return LoadContext(parent, ctx)
}
