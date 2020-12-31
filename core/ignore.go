package core

import (
	"io/ioutil"
	"strings"
)

// IgnoreRules contains default ignore rules.
// Use init func to append additional rules.
var IgnoreRules = []string{".git"}

// IgnoreFile is the name of ignore files.
const IgnoreFile = ".multignore"

// Ignore returns a list of files to ignore.
// If an ignore file exists its rules will
// be appended to the list of default rules.
func Ignore() ([]string, error) {
	if _, err := fs.Stat(IgnoreFile); err != nil {
		return IgnoreRules, nil
	}

	file, err := fs.Open(IgnoreFile)
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
