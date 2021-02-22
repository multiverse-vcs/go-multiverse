package command

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
	ignore "github.com/sabhiram/go-gitignore"
)

const (
	// DotDir is the name of the dot directory.
	DotDir = ".multi"
	// IgnoreFile is the name of the ignore file.
	IgnoreFile = "multi.ignore"
)

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{".git", "node_modules", DotDir}

// Context contains command context.
type Context struct {
	// Config contains repository settings.
	Config *Config
	// DAG contains all versioned files.
	DAG ipld.DAGService
	// Root is the top level directory.
	Root string
}

// NewContext returns a new context.
func NewContext(cwd string) (*Context, error) {
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

// Ignore returns ignore rules for the current context.
func (c *Context) Ignore() (*ignore.GitIgnore, error) {
	path := filepath.Join(c.Root, IgnoreFile)

	_, err := os.Lstat(path)
	if err == nil {
		return ignore.CompileIgnoreFileAndLines(path, IgnoreRules...)
	}

	if os.IsNotExist(err) {
		return ignore.CompileIgnoreLines(IgnoreRules...), nil
	}

	return nil, err
}
