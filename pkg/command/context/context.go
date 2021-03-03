package context

import (
	"errors"
	"os"
	"path/filepath"

	blockservice "github.com/ipfs/go-blockservice"
	badger "github.com/ipfs/go-ds-badger2"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"

	"github.com/multiverse-vcs/go-multiverse/internal/ignore"
)

// DotDir is the name of the dot directory.
const DotDir = ".multi"

// DefaultIgnore contans the default ignore rules.
var DefaultIgnore = ignore.New("", ".git", ".svn", ".hg", ".multi")

// Context contains command context.
type Context struct {
	// Blocks is the ipfs blockstore.
	Blocks blockstore.Blockstore
	// Config contains repository settings.
	Config *Config
	// DAG contains all versioned files.
	DAG ipld.DAGService
	// Root is the top level directory.
	Root string
}

// Init initializes a new context.
func Init(cwd string) error {
	if _, err := Root(cwd); err == nil {
		return errors.New("repo already exists")
	}

	root := filepath.Join(cwd, DotDir)
	if err := os.Mkdir(root, 0755); err != nil {
		return err
	}

	config := NewConfig(root)
	return config.Write()
}

// New returns a new context.
func New(cwd string) (*Context, error) {
	root, err := Root(cwd)
	if err != nil {
		return nil, err
	}

	config := NewConfig(root)
	if err := config.Read(); err != nil {
		return nil, err
	}

	dpath := filepath.Join(root, "datastore")
	dopts := badger.DefaultOptions

	dstore, err := badger.NewDatastore(dpath, &dopts)
	if err != nil {
		return nil, err
	}

	bstore := blockstore.NewBlockstore(dstore)
	exc := offline.Exchange(bstore)
	bserv := blockservice.New(bstore, exc)

	return &Context{
		Blocks: bstore,
		Config: config,
		DAG:    merkledag.NewDAGService(bserv),
		Root:   filepath.Dir(root),
	}, nil
}

// Root searches for the repository root.
func Root(root string) (string, error) {
	path := filepath.Join(root, DotDir)

	_, err := os.Lstat(path)
	if err == nil {
		return path, nil
	}

	if !os.IsNotExist(err) {
		return "", err
	}

	parent := filepath.Dir(root)
	if parent == root {
		return "", errors.New("repo not found")
	}

	return Root(parent)
}
