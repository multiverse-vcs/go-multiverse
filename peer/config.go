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
	// PrivateKey is the base64 encoded private key.
	PrivateKey string `json:"private_key"`

	path string
}

// GenerateConfig creates and saves a config with default settings.
func GenerateConfig(root string) (*Config, error) {
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

	if err := config.Save(); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadConfig returns a config from the given root dir.
// If it does not exist a new default config is returned.
func LoadConfig(root string) (*Config, error) {
	path := filepath.Join(root, ConfigFile)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return GenerateConfig(root)
	}

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

// Save writes the config to the path.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0644)
}
