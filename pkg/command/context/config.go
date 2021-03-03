package context

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

// Branch contains branch info.
type Branch struct {
	// Head is the CID of the branch head.
	Head cid.Cid `json:"head"`
	// Stash is the CID of the tree stash.
	Stash cid.Cid `json:"stash"`
	// Remote is the remote branch path.
	Remote string `json:"remote"`
}

// Config contains repository info.
type Config struct {
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
	// Branches contains named branches.
	Branches map[string]*Branch `json:"branches"`
	// Remotes contains named remotes.
	Remotes map[string]string `json:"remotes"`

	path string
}

// New returns a config with default settings.
func NewConfig(root string) *Config {
	return &Config{
		Branch: DefaultBranch,
		Branches: map[string]*Branch{
			DefaultBranch: {},
		},
		Remotes: make(map[string]string),
		path:    filepath.Join(root, ConfigFile),
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
