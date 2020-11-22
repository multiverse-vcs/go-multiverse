package core

import (
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/sabhiram/go-gitignore"
)

// Worktree adds the current working tree to the merkle dag.
// Optional ignore rules can be used to filter out files.
func (c *Context) Worktree() (ipld.Node, error) {
	rules, err := c.Ignore()
	if err != nil {
		return nil, err
	}

	filter, err := ignore.CompileIgnoreLines(rules...)
	if err != nil {
		return nil, err
	}

	return c.Add("", filter)
}
