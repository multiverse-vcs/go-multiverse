package ignore

import (
	"path/filepath"
	"strings"
)

// Rule is used to match file paths.
type Rule interface {
	// Match returns true if the path matches.
	Match(string) (bool, error)
}

// compile time interface checks
var _ Rule = (*NegateRule)(nil)
var _ Rule = (BaseRule)("")
var _ Rule = (PathRule)("")
var _ Rule = (EmptyRule)("")

// NegateRule negates a rule.
type NegateRule struct {
	wrapped Rule
}

// BaseRule matches the base of the path.
type BaseRule string

// PathRule matches the entire path.
type PathRule string

// EmptyRule never matches.
type EmptyRule string

// Match returns true if the path matches.
func (r NegateRule) Match(name string) (bool, error) {
	match, err := r.wrapped.Match(name)
	return !match, err
}

// Match returns true if the path matches.
func (r BaseRule) Match(name string) (bool, error) {
	base := filepath.Base(name)
	return filepath.Match(string(r), base)
}

// Match returns true if the path matches.
func (r PathRule) Match(name string) (bool, error) {
	return filepath.Match(string(r), name)
}

// Match returns true if the path matches.
func (r EmptyRule) Match(name string) (bool, error) {
	return false, nil
}

// ParseRule returns a rule for the given pattern.
func ParseRule(dir, pattern string) Rule {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return EmptyRule(pattern)
	}

	if strings.HasPrefix(pattern, "#") {
		return EmptyRule(pattern)
	}

	var rule Rule
	if strings.Contains(pattern, "/") {
		rule = PathRule(filepath.Join(dir, pattern))
	} else {
		rule = BaseRule(pattern)
	}

	if strings.HasPrefix(pattern, "!") {
		rule = NegateRule{rule}
	}

	return rule
}
