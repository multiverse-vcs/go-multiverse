package command

import (
	"encoding/json"
	"os"
	"path/filepath"

	cid "github.com/ipfs/go-cid"
)

const (
	// ConfigFile is the name of the config file.
	ConfigFile = "config.json"
	// DefaultBranch is the name of the default branch.
	DefaultBranch = "main"
)

// Config contains repository info.
type Config struct {
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
	// Branches is a map of names to commit CIDs.
	Branches map[string]cid.Cid `json:"branches"`
	// Index is the CID of the current commit.
	Index cid.Cid `json:"index"`
	// Remote is the repository remote server.
	Remote string `json:"remote"`
	// Tags is a map of names to commit CIDs.
	Tags map[string]cid.Cid `json:"tags"`

	path string
}

// New returns a config with default settings.
func NewConfig(root string) *Config {
	return &Config{
		Branch:   DefaultBranch,
		Branches: make(map[string]cid.Cid),
		Tags:     make(map[string]cid.Cid),
		path:     filepath.Join(root, ConfigFile),
	}
}

// Read reads the config from the path.
func (c *Config) Read() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

// Write writes the config to the path.
func (c *Config) Write() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0644)
}
