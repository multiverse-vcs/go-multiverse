package peer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/p2p"
)

// ConfigFile is the name of the config file.
const ConfigFile = "config.json"

// Config contains peer settings.
type Config struct {
	// Author contains the current author.
	Author *data.Author `json:"author"`
	// Sequence is the version number to publish.
	Sequence uint64 `json:"sequence"`
	// PrivateKey is the base64 encoded private key.
	PrivateKey string `json:"private_key"`

	path string
}

// NewConfig creates a config with default settings.
func NewConfig(root string) (*Config, error) {
	priv, err := p2p.GenerateKey()
	if err != nil {
		return nil, err
	}

	encoded, err := p2p.EncodeKey(priv)
	if err != nil {
		return nil, err
	}

	config := Config{
		Author:     data.NewAuthor(),
		PrivateKey: encoded,
		path:       filepath.Join(root, ConfigFile),
	}

	return &config, nil
}

// LoadConfig returns a config from the given root dir.
func LoadConfig(root string) (*Config, error) {
	path := filepath.Join(root, ConfigFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.path = path
	return &config, nil
}

// OpenConfig loads the config from the given root dir.
// If the config does not exist a new one is generated.
func OpenConfig(root string) (*Config, error) {
	config, err := LoadConfig(root)
	if err == nil {
		return config, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	config, err = NewConfig(root)
	if err != nil {
		return nil, err
	}

	return config, config.Save()
}

// Save writes the config to the path.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0644)
}
