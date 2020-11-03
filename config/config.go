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

// DefaultConfig is the name of the default repo config file.
const DefaultConfig = ".multiverse.json"

var (
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo cannot be found.
	ErrRepoNotFound = errors.New("repo not found")
)

// Config contains local repo info.
type Config struct {
	// Path is the repo root directory.
	Path string `json:"-"`
	// Head is the CID of latest commit.
	Head cid.Cid `json:"head"`
}

// Init creates a new config at the given path.
func Init(path string) (*Config, error) {
	_, err := Open(path)
	if err == nil {
		return nil, ErrRepoExists
	}

	c := Config{Path: path}
	if err := c.Write(); err != nil {
		return nil, err
	}

	return &c, nil
}

// Open reads a config in the path or parent directories.
func Open(path string) (*Config, error) {
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

// Read reads a config in the current directory.
func Read(path string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, DefaultConfig))
	if err != nil {
		return nil, err
	}

	c := Config{Path: path}
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Write writes the config to the root directory.
func (c *Config) Write() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	name := filepath.Join(c.Path, DefaultConfig)
	return ioutil.WriteFile(name, data, 0644)
}
