package core

import (
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/sabhiram/go-gitignore"
)

// IgnoreFile is the name of ignore files.
const IgnoreFile = ".multignore"

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{".multiverse"}

// Worktree adds the current working tree to the merkle dag.
func (c *Context) Worktree() (ipld.Node, error) {
	filter, err := ignore.CompileIgnoreLines(IgnoreRules...)
	if err != nil {
		return nil, err
	}

	return c.Add(c.config.Root, filter)
}
