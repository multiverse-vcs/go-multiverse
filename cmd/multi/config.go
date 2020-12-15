package main

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/spf13/afero"
)

const (
	// ConfigFile is the name of the config file.
	ConfigFile = "config"
	// KeyType is the default private key type.
	KeyType = crypto.Ed25519
)

// Branch contains local branch info.
type Branch struct {
	// Head is the tip of the branch.
	Head cid.Cid
	// Stash is a copy of the working directory.
	Stash cid.Cid
}

// Config contains local repo info.
type Config struct {
	// Index is the commit changes are based on.
	Index cid.Cid
	// PrivateKey is the key used for p2p networking.
	PrivateKey string
	// Branch is the name of the current branch.
	Branch string
	// Branches contains a map of local branches.
	Branches map[string]*Branch
	// Remotes contains a map of remote endpoints.
	Remotes map[string]string
}

// DefaultConfig returns a config with default values.
func DefaultConfig() (*Config, error) {
	cfg := &Config{
		Branch: "default",
		Branches: map[string]*Branch{
			"default": {},
		},
		Remotes: map[string]string{
			"local": "http://127.0.0.1:5001",
		},
	}

	return cfg, cfg.GenerateKey()
}

// ReadConfig reads the config from the given fs.
func ReadConfig(root string, cfg *Config) error {
	path := filepath.Join(root, ConfigFile)
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg)
}

// WriteConfig writes the config to the given fs.
func WriteConfig(root string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	path := filepath.Join(root, ConfigFile)
	return afero.WriteFile(fs, path, data, 0644)
}

// Head returns the current branch head.
func (c *Config) Head() cid.Cid {
	branch, err := c.GetBranch(c.Branch)
	if err != nil {
		return cid.Cid{}
	}

	return branch.Head
}

// SetHead sets the current branch head.
func (c *Config) SetHead(head cid.Cid) {
	branch, err := c.GetBranch(c.Branch)
	if err != nil {
		return
	}

	branch.Head = head
}

// Stash returns the current branch stash.
func (c *Config) Stash() cid.Cid {
	branch, err := c.GetBranch(c.Branch)
	if err != nil {
		return cid.Cid{}
	}

	return branch.Stash
}

// SetStash sets the current branch stash.
func (c *Config) SetStash(stash cid.Cid) {
	branch, err := c.GetBranch(c.Branch)
	if err != nil {
		return
	}

	branch.Stash = stash
}

// Key returns the private key from the config.
func (c *Config) Key() (crypto.PrivKey, error) {
	data, err := crypto.ConfigDecodeKey(c.PrivateKey)
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(data)
}

// SetKey sets the config private key.
func (c *Config) SetKey(key crypto.PrivKey) error {
	data, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return err
	}

	c.PrivateKey = crypto.ConfigEncodeKey(data)
	return nil
}

// GenerateKey creates and sets a new private key.
func (c *Config) GenerateKey() error {
	priv, _, err := crypto.GenerateKeyPair(KeyType, -1)
	if err != nil {
		return err
	}

	return c.SetKey(priv)
}

// AddBranch creates a new branch with the given name and head.
func (c *Config) AddBranch(name string, head cid.Cid) error {
	if _, ok := c.Branches[name]; ok {
		return errors.New("branch already exists")
	}

	c.Branches[name] = &Branch{Head: head}
	return nil
}

// DeleteBranch deletes the branch with the given name.
func (c *Config) DeleteBranch(name string) error {
	if _, ok := c.Branches[c.Branch]; !ok {
		return errors.New("branch does not exist")
	}

	delete(c.Branches, name)
	return nil
}

// GetBranch returns the branch with the given name
func (c *Config) GetBranch(name string) (*Branch, error) {
	branch, ok := c.Branches[c.Branch]
	if !ok {
		return nil, errors.New("branch does not exist")
	}

	return branch, nil
}
