// Package repo contains methods for managing repositories.
package repo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/core"
)

// Config is the name of the repo config file.
const Config = ".multiverse"

// Repo contains project info.
type Repo struct {
	// Path is the path to the repo config.
	Path string `json:"-"`
	// Root is the path to the root directory.
	Root string `json:"-"`
	// Name is the human friendly name of the repo.
	Name string `json:"name"`
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
	// Branches is a map of branch heads.
	Branches map[string]string `json:"branches"`
}

func init() {
	core.IgnoreRules = append(core.IgnoreRules, Config)
}

// Default returns a config with default settings.
func Default(root string, name string) *Repo {
	return &Repo{
		Name:     name,
		Branch:   "default",
		Branches: make(map[string]string),
		Path:     filepath.Join(root, Config),
		Root:     root,
	}
}

// Find searches for the config in parent directories.
func Find(root string) (string, error) {
	path := filepath.Join(root, Config)

	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return path, nil
	}

	parent := filepath.Dir(root)
	if parent == root {
		return "", errors.New("repo not found")
	}

	return Find(parent)
}

// Read reads the repo from the given path.
func Read(root string) (*Repo, error) {
	path, err := Find(root)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var repo Repo
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}

	repo.Path = path
	repo.Root = filepath.Dir(path)

	return &repo, nil
}

// SetHead sets the cid of the current repo branch.
func (repo *Repo) SetHead(id cid.Cid) {
	repo.Branches[repo.Branch] = id.String()
}

// Head returns the cid of the current repo branch.
func (repo *Repo) Head() (cid.Cid, error) {
	id, ok := repo.Branches[repo.Branch]
	if !ok {
		return cid.Cid{}, nil
	}

	return cid.Parse(id)
}

// Write saves the repo config.
func (repo *Repo) Write() error {
	data, err := json.MarshalIndent(repo, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(repo.Path, data, 0644)
}
