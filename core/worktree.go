package core

import (
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
)

// IgnoreFile is the name of ignore files.
const IgnoreFile = ".multignore"

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{".multiverse"}

// Worktree returns the current working directory.
func (c *Context) Worktree() (files.Node, error) {
	info, err := os.Stat(c.root)
	if err != nil {
		return nil, err
	}

	ignore := filepath.Join(c.root, IgnoreFile)
	if _, err := os.Stat(ignore); err != nil {
		ignore = ""
	}

	filter, err := files.NewFilter(ignore, IgnoreRules, true)
	if err != nil {
		return nil, err
	}

	return files.NewSerialFileWithFilter(c.root, filter, info)
}
