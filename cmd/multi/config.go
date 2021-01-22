package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-cid"
)

const (
	// ConfigFile is the name of the config file.
	ConfigFile = "multi.json"
	// IgnoreFile is the name of the ignore file.
	IgnoreFile = "multi.ignore"
)

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{".git", ConfigFile}

// Config contains project info.
type Config struct {
	// Path is the path to the config file.
	Path string `json:"-"`
	// Root is the path to the root directory.
	Root string `json:"-"`
	// Repo is the CID of the repository.
	Repo cid.Cid `json:"repo"`
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
	// Index is the CID of the current commit.
	Index cid.Cid `json:"index"`
}

// NewConfig returns a config with default settings.
func NewConfig(root string) *Config {
	return &Config{
		Branch: "default",
		Path:   filepath.Join(root, ConfigFile),
		Root:   root,
	}
}

// FindConfig searches for the config in parent directories.
func FindConfig(root string) (string, error) {
	path := filepath.Join(root, ConfigFile)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	parent := filepath.Dir(root)
	if parent == root {
		return "", errors.New("repo not found")
	}

	return FindConfig(parent)
}

// LoadConfig reads the repo from the given path.
func LoadConfig(root string) (*Config, error) {
	path, err := FindConfig(root)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.Path = path
	config.Root = filepath.Dir(path)

	return &config, nil
}

// Ignore returns a list of files to ignore.
func (c *Config) Ignore() ([]string, error) {
	path := filepath.Join(c.Root, IgnoreFile)
	if _, err := os.Stat(path); err != nil {
		return IgnoreRules, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return append(IgnoreRules, lines...), nil
}

// Write saves the repo config.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Path, data, 0644)
}
