// Package config contains methods for reading and writing configurations.
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
)

const (
	// ConfigFile is the name of the config file.
	ConfigFile = ".multiverse.json"
	// DefaultBranch is the name of the default branch.
	DefaultBranch = "default"
)

var (
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo does not exist.
	ErrRepoNotFound = errors.New("repo not found")
	// ErrBranchNotFound is returned when a branch does not exist.
	ErrBranchNotFound = errors.New("branch does not exist")
	// ErrBranchDetached is returned when base is behind head.
	ErrBranchDetached = errors.New("branch is detached")
)

// Config contains local repo info.
type Config struct {
	// Path is the repo root directory.
	Path string `json:"-"`
	// Base is the cid of the current working base.
	Base cid.Cid `json:"base"`
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
	// Branches is a map of branch heads.
	Branches map[string]cid.Cid `json:"branches"`
}

// Init creates a new config at the given path.
func Init(path string) (*Config, error) {
	_, err := Open(path)
	if err == nil {
		return nil, ErrRepoExists
	}

	c := Config{Branch: DefaultBranch}
	if err := c.Write(); err != nil {
		return nil, err
	}

	c.Path = path
	return &c, nil
}

// Open searches for a config in the path or parent directories.
func Open(path string) (*Config, error) {
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

// Read reads a config in the current directory.
func Read(path string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, ConfigFile))
	if err != nil {
		return nil, err
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	c.Path = path
	return &c, nil
}

// Detached returns an error when base is not equal to head.
func (c *Config) Detached() error {
	head, err := c.Head()
	if err != nil {
		return err
	}

	if head != c.Base {
		return ErrBranchDetached
	}

	return nil
}

// Head returns the tip of the current branch.
func (c *Config) Head() (cid.Cid, error) {
	head, ok := c.Branches[c.Branch]
	if !ok {
		return cid.Cid{}, ErrBranchNotFound
	}

	return head, nil
}

// Write writes the config to the root directory.
func (c *Config) Write() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(c.Path, ConfigFile), data, 0644)
}
