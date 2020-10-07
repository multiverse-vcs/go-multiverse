package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
)

const (
	// DefaultConfig is the name of the default repo config file.
	DefaultConfig = "multi.json"
	// DefaultBranch is the name of the default repo branch.
	DefaultBranch = "default"
)

var (
	// ErrRepoExists is returned when a repo already exists.
	ErrRepoExists = errors.New("repo already exists")
	// ErrRepoNotFound is returned when a repo cannot be found.
	ErrRepoNotFound = errors.New("repo not found")
)

// Config contains local repo info.
type Config struct {
	// Path repo root directory.
	Path string `json:"-"`
	// Head is the CID of latest commit.
	Head cid.Cid `json:"head"`
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
}

// InitConfig creates a new config at the given path.
func InitConfig(path string, head cid.Cid) (*Config, error) {
	_, err := OpenConfig(path)
	if err == nil {
		return nil, ErrRepoExists
	}

	c := Config{path, head, DefaultBranch}
	if err := c.Write(); err != nil {
		return nil, err
	}

	return &c, nil
}

// OpenConfig reads a config in the current or parent directories.
func OpenConfig(path string) (*Config, error) {
	_, err := os.Stat(filepath.Join(path, DefaultConfig))
	if err == nil {
		return ReadConfig(path)
	}

	parent := filepath.Dir(path)
	if parent == path {
		return nil, ErrRepoNotFound
	}

	return OpenConfig(parent)
}

// ReadConfig reads a config in the current directory.
func ReadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, DefaultConfig))
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

// Write writes the config to the root directory.
func (c *Config) Write() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(DefaultConfig, data, 0644)
}
