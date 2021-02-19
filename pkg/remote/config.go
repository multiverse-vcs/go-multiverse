package remote

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// ConfigFile is the name of the config file.
const ConfigFile = "config.json"

// Config contains repository info.
type Config struct {
	// Author contains published repositories.
	Author *object.Author
	// PrivateKey is the private key of the remote.
	PrivateKey string

	path string
}

// New returns a config with default settings.
func NewConfig(root string) *Config {
	return &Config{
		Author: object.NewAuthor(),
		path:   filepath.Join(root, ConfigFile),
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
