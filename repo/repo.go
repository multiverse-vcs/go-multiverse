// Package repo contains methods for working with local repos.
package repo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
)

const (
	// ConfigFile is the name of the config file.
	ConfigFile = ".multiverse.json"
	// DefaultBranch is the name of the default branch.
	DefaultBranch = "default"
	// IgnoreFile is the name of the ignore file.
	IgnoreFile = ".multiverse.ignore"
)

var (
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo does not exist.
	ErrRepoNotFound = errors.New("repo not found")
	// ErrRepoDetached is returned when base is behind head.
	ErrRepoDetached = errors.New("repo base is behind head")
)

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{ConfigFile}

// Repo contains local repo configuration.
type Repo struct {
	// Path is the repo root directory.
	Path string `json:"-"`
	// Base is the cid of the current working base.
	Base cid.Cid `json:"base"`
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
	// Branches is a map of branch heads.
	Branches Branches `json:"branches"`
}

// NewRepo returns a new repo with default values.
func NewRepo(path string) *Repo {
	return &Repo{
		Path: path, 
		Branch: DefaultBranch,
		Branches: Branches{},
	}
}

// Init creates a new repo at the given path.
func Init(path string) (*Repo, error) {
	_, err := Open(path)
	if err == nil {
		return nil, ErrRepoExists
	}

	r := NewRepo(path)
	if err := r.Write(); err != nil {
		return nil, err
	}

	return r, nil
}

// Open searches for a repo in the path or parent directories.
func Open(path string) (*Repo, error) {
	_, err := os.Stat(filepath.Join(path, ConfigFile))
	if err == nil {
		return Read(path)
	}

	parent := filepath.Dir(path)
	if parent == path {
		return nil, ErrRepoNotFound
	}

	return Open(parent)
}

// Read reads a repo in the current directory.
func Read(path string) (*Repo, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, ConfigFile))
	if err != nil {
		return nil, err
	}

	r := Repo{Path: path}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// Detached returns an error when base is not equal to head.
func (r *Repo) Detached() error {
	head, err := r.Branches.Head(r.Branch)
	if err != nil {
		return err
	}

	if head != r.Base {
		return ErrRepoDetached
	}

	return nil
}

// Ignore returns a files.Filter that is used to ignore files.
func (r *Repo) Ignore() (*files.Filter, error) {
	ignore := filepath.Join(r.Path, IgnoreFile)
	if _, err := os.Stat(ignore); err != nil {
		ignore = ""
	}

	return files.NewFilter(ignore, IgnoreRules, true)
}

// Tree returns the repo working tree files.Node.
func (r *Repo) Tree() (files.Node, error) {
	info, err := os.Stat(r.Path)
	if err != nil {
		return nil, err
	}

	ignore, err := r.Ignore()
	if err != nil {
		return nil, err
	}

	return files.NewSerialFileWithFilter(r.Path, ignore, info)
}

// Write writes the config to the root directory.
func (r *Repo) Write() error {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(r.Path, ConfigFile), data, 0644)
}
