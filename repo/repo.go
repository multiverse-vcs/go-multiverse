package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/multiformats/go-multihash"
	"github.com/spf13/afero"
)

// Contains repo info
type Repo struct {
	fs   afero.Afero
	ipfs iface.CoreAPI
}

// Create a repo
func NewRepo() *Repo {
	ipfs, err := httpapi.NewLocalApi()
	if err != nil {
		panic(err)
	}

	fs := afero.Afero{Fs: afero.NewOsFs()}
	return &Repo{fs: fs, ipfs: ipfs}
}

// Initialize a repo in the directory.
func (repo *Repo) Init(path string) error {
	dir := filepath.Join(path, ".multi")
	return repo.fs.MkdirAll(dir, os.ModePerm)
}

// Walk parent directories until repo root is found.
func (repo *Repo) Root() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	multi := filepath.Join(cwd, ".multi")
	if exists, _ := repo.fs.DirExists(multi); exists {
		return nil
	}

	parent, err := filepath.Abs("..")
	if err != nil {
		return err
	}

	if parent == cwd {
		return fmt.Errorf("Repo not found")
	}

	if err := os.Chdir(parent); err != nil {
		return err
	}

	return repo.Root()
}

// Return the repo root directory.
func (repo *Repo) Dir(path string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	defer os.Chdir(cwd)

	if err := os.Chdir(path); err != nil {
		return "", err
	}

	if err := repo.Root(); err != nil {
		return "", err
	}

	return os.Getwd()
}

// Record changes in the local repo.
func (repo *Repo) Commit(message string) error {
	if err := repo.Root(); err != nil {
		return err
	}

	stat, err := repo.fs.Stat(".")
	if err != nil {
		return err
	}

	filter, err := files.NewFilter("", []string{".multi"}, true)
	if err != nil {
		return err
	}

	node, err := files.NewSerialFileWithFilter(".", filter, stat)
	if err != nil {
		return err
	}

	changes, err := repo.ipfs.Unixfs().Add(context.TODO(), node)
	if err != nil {
		return err
	}

	commit := make(map[string]interface{})
	commit["message"] = message
	commit["changes"] = changes.Cid()

	dag, err := cbornode.WrapObject(commit, multihash.SHA2_256, -1)
	if err != nil {
		return err
	}

	return repo.ipfs.Dag().Pinning().Add(context.TODO(), dag)
}