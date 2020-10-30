package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
)

// DefaultConfig is the name of the default repo config file.
const DefaultConfig = ".multiverse.json"

// Config contains local repo info.
type Config struct {
	// Path repo root directory.
	Path string `json:"-"`
	// Head is the CID of latest commit.
	Head cid.Cid `json:"head"`
}

// InitConfig creates a new config at the given path.
func InitConfig(path string) (*Config, error) {
	_, err := OpenConfig(path)
	if err == nil {
		return nil, ErrRepoExists
	}

	c := Config{Path: path}
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
