package core

import (
	"io/ioutil"
	"strings"
)

// IgnoreRules contains default ignore rules.
// Use init func to append additional rules.
var IgnoreRules = []string{}

// IgnoreFile is the name of ignore files.
const IgnoreFile = ".multignore"

// Ignore returns a list of files to ignore.
// If an ignore file exists its rules will
// be appended to the list of default rules.
func (c *Context) Ignore() ([]string, error) {
	path := c.Fs.Join(c.Fs.Root(), IgnoreFile)
	if _, err := c.Fs.Lstat(path); err != nil {
		return IgnoreRules, nil
	}

	file, err := c.Fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	lines = append(IgnoreRules, lines...)

	return lines, nil
}
