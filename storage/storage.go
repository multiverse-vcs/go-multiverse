// Package storage implements storage strategies.
package storage

import (
	"encoding/json"

	"github.com/ipfs/go-ipfs-blockstore"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/spf13/afero"
)

const (
	// DotDir is the name of the dot directory.
	DotDir = ".multiverse"
	// DataDir is the name of the data directory.
	DataDir = "datastore"
	// ConfigFile is the name of the config file.
	ConfigFile = "config.json"
)

// Store contains storage services.
type Store struct {
	// Dag is the merkledag interface.
	Dag ipld.DAGService
	// Dot is the multiverse dot directory.
	Dot afero.Fs
	// Cwd is the current working directory.
	Cwd afero.Fs
	// Host is the optional p2p host.
	Host host.Host
	// Router is the optional p2p routing service.
	Router routing.Routing

	bstore blockstore.Blockstore
}

// ReadConfig reads the config file from the given filesystem.
func (s *Store) ReadConfig() (*config.Config, error) {
	data, err := afero.ReadFile(s.Dot, ConfigFile)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// WriteConfig writes the config file to the given filesystem.
func (s *Store) WriteConfig(c *config.Config) error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return afero.WriteFile(s.Dot, ConfigFile, data, 0644)
}
