package core

import (
	"io/ioutil"
	"strings"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/sabhiram/go-gitignore"
)

// IgnoreFile is the name of ignore files.
const IgnoreFile = ".multignore"

// IgnoreRules contains default ignore rules.
var IgnoreRules = []string{".multiverse"}

// Worktree adds the current working tree to the merkle dag.
func (c *Context) Worktree() (ipld.Node, error) {
	filter, err := c.ignore()
	if err != nil {
		return nil, err
	}

	return c.Add(c.fs.Root(), filter)
}

func (c *Context) ignore() (*ignore.GitIgnore, error) {
	path := c.fs.Join(c.fs.Root(), IgnoreFile)
	if _, err := c.fs.Lstat(path); err != nil {
		return ignore.CompileIgnoreLines(IgnoreRules...)
	}

	file, err := c.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	lines = append(lines, IgnoreRules...)

	return ignore.CompileIgnoreLines(lines...)
}
