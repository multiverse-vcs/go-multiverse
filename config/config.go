// Package config contains configuration definitions.
package config

import (
	"errors"

	"github.com/ipfs/go-cid"
)

const (
	DefaultBranch    = "default"
	DefaultRemote    = "local"
	DefaultRemoteURL = "http://127.0.0.1:5001"
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
	// Branch is the name of the current branch.
	Branch string
	// Branches contains a map of local branches.
	Branches map[string]*Branch
	// Remotes contains a map of remote endpoints.
	Remotes map[string]string
}

// Default returns a new config with default settings.
func Default() *Config {
	return &Config{
		Branch: DefaultBranch,
		Branches: map[string]*Branch{
			DefaultBranch: {},
		},
		Remotes: map[string]string{
			DefaultRemote: DefaultRemoteURL,
		},
	}
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
