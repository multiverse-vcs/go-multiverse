// Package config contains configuration definitions.
package config

import (
	"encoding/json"
	"errors"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	fsutil "github.com/go-git/go-billy/v5/util"
	"github.com/ipfs/go-cid"
)

// DefaultBranch is the name of the default branch.
const DefaultBranch = "default"

// ConfigName is the name of the config file.
const ConfigName = "config.json"

// Config contains configuration info.
type Config struct {
	Base     cid.Cid            `json:"base"`
	Branch   string             `json:"branch"`
	Branches map[string]cid.Cid `json:"branches"`
	Head     cid.Cid            `json:"head"`

	fs billy.Filesystem
}

// NewMockConfig returns a config that can be used for testing.
func NewMockConfig() *Config {
	return &Config{
		Branch: DefaultBranch,
		fs:     memfs.New(),
	}
}

// Detached returns an error if base is not equal to head.
func (c *Config) Detached() error {
	if c.Base != c.Head {
		return errors.New("base is behind head")
	}

	return nil
}

// Write persists the config to the filesystem.
func (c *Config) Write() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return fsutil.WriteFile(c.fs, ConfigName, data, 0644)
}
