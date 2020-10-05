package repo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
)

// DefaultConfig is the name of the default repo config file.
const DefaultConfig = "multi.json"

var (
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo cannot be found.
	ErrRepoNotFound = errors.New("repo not found")
)

// Repo contains repo info.
type Repo struct {
	// Path repo root directory.
	Path string `json:"-"`
	// Head is the CID of latest commit.
	Head cid.Cid `json:"head"`
}

// Init creates a new empty repo at the given path.
func Init(path string) (*Repo, error) {
	_, err := Open(path)
	if err == nil {
		return nil, ErrRepoExists
	}

	r := Repo{Path: path}
	return &r, r.Write()
}

// Open returns an existing repo in the current or parent directories.
func Open(path string) (*Repo, error) {
	_, err := os.Stat(filepath.Join(path, DefaultConfig))
	if err == nil {
		return Read(path)
	}

	parent := filepath.Dir(path)
	if parent == path {
		return nil, ErrRepoNotFound
	}

	return Open(parent)
}

// Read returns an existing repo in the current directory.
func Read(path string) (*Repo, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, DefaultConfig))
	if err != nil {
		return nil, err
	}

	r := Repo{Path: path}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// Write saves the repo config to the root directory.
func (r *Repo) Write() error {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(DefaultConfig, data, 0644)
}
