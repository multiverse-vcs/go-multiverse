package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/go-git/go-billy/v5"
	fsutil "github.com/go-git/go-billy/v5/util"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/core"
)

const (
	// DotDir is the name of the dot directory.
	DotDir = ".multiverse"
	// DatastoreName is the name of the datastore directory.
	DatastoreDir = "datastore"
	// ConfigFile is the name of the config file.
	ConfigFile = "config.json"
)

func init() {
	core.IgnoreRules = append(core.IgnoreRules, DotDir, ".git")
}

// Context is the context for cli commands.
type Context struct {
	core.Context

	dot billy.Filesystem
}

// NewContext is used internally to create a command context.
func NewContext(fs billy.Filesystem, ctx context.Context) (*Context, error) {
	dot, err := fs.Chroot(DotDir)
	if err != nil {
		return nil, err
	}

	cfg, err := ReadConfig(dot)
	if err != nil {
		return nil, err
	}

	dstore, err := NewDatastore(dot)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(dstore)
	bserv := blockservice.New(bstore, offline.Exchange(bstore))
	dserv := merkledag.NewDAGService(bserv)

	corectx := core.Context{
		Context: ctx,
		Config:  cfg,
		Dag:     dserv,
		Fs:      fs,
	}

	return &Context{
		Context: corectx,
		dot:     dot,
	}, nil
}

// InitContext initializes a context in the given path.
func InitContext(fs billy.Filesystem, ctx context.Context) error {
	if _, err := LoadContext(fs, ctx); err == nil {
		return errors.New("repo already exists")
	}

	if err := fs.MkdirAll(DotDir, 0755); err != nil {
		return err
	}

	dot, err := fs.Chroot(DotDir)
	if err != nil {
		return err
	}

	cfg := config.Config{
		Branch: config.DefaultBranch,
	}

	return WriteConfig(dot, &cfg)
}

// LoadContext loads a context in the given path or parent directories.
func LoadContext(fs billy.Filesystem, ctx context.Context) (*Context, error) {
	info, err := fs.Lstat(DotDir)
	if err == nil && info.IsDir() {
		return NewContext(fs, ctx)
	}

	parent, err := fs.Chroot("..")
	if err != nil {
		return nil, err
	}

	if parent.Root() == fs.Root() {
		return nil, errors.New("repo not found")
	}

	return LoadContext(parent, ctx)
}

// NewDatastore returns a datastore backed by the given filesystem.
func NewDatastore(fs billy.Filesystem) (datastore.Batching, error) {
	type fsBased interface {
		Filesystem() billy.Filesystem
	}

	opts := badger.DefaultOptions
	if _, ok := fs.(fsBased); ok {
		opts.WithInMemory(true)
	}

	path := fs.Join(fs.Root(), DatastoreDir)
	return badger.NewDatastore(path, &opts)
}

// ReadConfig reads the config file from the given filesystem.
func ReadConfig(fs billy.Filesystem) (*config.Config, error) {
	file, err := fs.Open(ConfigFile)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// WriteConfig writes the config file to the given filesystem.
func WriteConfig(fs billy.Filesystem, c *config.Config) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return fsutil.WriteFile(fs, ConfigFile, data, 0644)
}
