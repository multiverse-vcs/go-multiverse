// Package config contains configuration definitions.
package config

import (
	"errors"

	"github.com/ipfs/go-cid"
)

// DefaultBranch is the name of the default branch.
const DefaultBranch = "default"

// Config contains local repo info.
type Config struct {
	// Index is the commit changes are based on.
	Index cid.Cid `json:"base"`
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
	// Branches contains a map of local branches.
	Branches map[string]cid.Cid
}

// Default returns a new config with default settings.
func Default() *Config {
	return &Config{
		Branch:   DefaultBranch,
		Branches: map[string]cid.Cid{DefaultBranch: {}},
	}
}

// Head returns the current branch head.
func (c *Config) Head() cid.Cid {
	head, _ := c.Branches[c.Branch]
	return head
}

// SetHead sets the current branch head to the given head.
func (c *Config) SetHead(head cid.Cid) {
	c.Branches[c.Branch] = head
}

// AddBranch creates a new branch with the given name and head.
func (c *Config) AddBranch(name string, head cid.Cid) error {
	if _, ok := c.Branches[name]; ok {
		return errors.New("branch already exists")
	}

	c.Branches[name] = head
	return nil
}

// DeleteBranch deletes the branch with the given name.
func (c *Config) DeleteBranch(name string) error {
	if _, ok := c.Branches[name]; !ok {
		return errors.New("branch does not exist")
	}

	delete(c.Branches, name)
	return nil
}
