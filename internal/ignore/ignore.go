package ignore

import (
	"os"
	"path/filepath"
	"strings"
)

// IgnoreFile is the name of the ignore file.
const IgnoreFile = ".multignore"

// Filter is a list of paths to ignore.
type Filter []Rule

// New returns a new ignore filter.
func New(dir string, patterns ...string) Filter {
	var rules []Rule
	for _, p := range patterns {
		rules = append(rules, ParseRule(dir, p))
	}

	return Filter(rules)
}

// Load returns the ignore filter from the given directory.
func Load(dir string) (Filter, error) {
	data, err := os.ReadFile(filepath.Join(dir, IgnoreFile))
	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	patterns := strings.Split(string(data), "\n")
	return New(dir, patterns...), nil
}

// Match returns true if the path matches any ignore rules.
func (f Filter) Match(name string) bool {
	for _, r := range f {
		if match, _ := r.Match(name); match {
			return true
		}
	}

	return false
}

// Merge combines the ignore rules with the other ignore.
func (f Filter) Merge(other Filter) Filter {
	var merge Filter
	for _, r := range other {
		merge = append(merge, r)
	}

	return append(f, merge...)
}
