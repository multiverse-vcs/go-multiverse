package command

import (
	"encoding/json"
	"os"
	"path/filepath"

	cid "github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
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
	// Index is the CID of the current commit.
	Index cid.Cid `json:"index"`
	// Remote is the repository remote server.
	Remote string `json:"remote"`
	// Repository contains references to branch heads.
	Repository *object.Repository `json:"repository"`

	path string
}

// New returns a config with default settings.
func NewConfig(root string) *Config {
	return &Config{
		Branch:     DefaultBranch,
		Repository: object.NewRepository(),
		path:       filepath.Join(root, ConfigFile),
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
