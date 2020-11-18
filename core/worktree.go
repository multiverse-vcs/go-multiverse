package core

import (
	"errors"

	ipld "github.com/ipfs/go-ipld-format"
)

// IgnoreFile is the name of ignore files.
const IgnoreFile = ".multignore"

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{".multiverse"}

// Worktree adds the current working tree to the merkle dag.
func (c *Context) Worktree() (ipld.Node, error) {
	info, err := c.fs.Lstat(c.config.Root)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("invalid worktree")
	}

	adder, err := c.NewAdder()
	if err != nil {
		return nil, err
	}

	// TODO implement ignore filters
	return adder.Add(c.config.Root, info)
}
