// Package core contains the core Multiverse types and functions.
package core

import (
	"errors"
)

var (
	// ErrNoLink is returned when an invalid path is resolved.
	ErrNoLink      = errors.New("no such link")
	// ErrZeroLenPath is returned when a path of len zero is resolved.
	ErrZeroLenPath = errors.New("zero length path")
	// ErrOutOfRange is returned when a path index is out of range.
	ErrOutOfRange  = errors.New("index out of range")
)