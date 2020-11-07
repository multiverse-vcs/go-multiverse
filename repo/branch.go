package repo

import (
	"errors"

	"github.com/ipfs/go-cid"
)

var (
	// ErrBranchExists is returned when a branch already exists.
	ErrBranchExists = errors.New("branch already exists")
	// ErrBranchNotFound is returned when a branch does not exist.
	ErrBranchNotFound = errors.New("branch does not exist")
)

// Branches is a map of names to cids.
type Branches map[string]cid.Cid

// Add appends a new branch with the given name and head.
func (b Branches) Add(name string, head cid.Cid) error {
	if _, ok := b[name]; ok {
		return ErrBranchExists
	}

	b[name] = head
	return nil
}

// Remove deletes the branch with the given name.
func (b Branches) Remove(name string) error {
	if _, ok := b[name]; !ok {
		return ErrBranchNotFound
	}

	delete(b, name)
	return nil
}

// Head returns the tip of the branch with the given name.
func (b Branches) Head(name string) (cid.Cid, error) {
	head, ok := b[name]
	if !ok {
		return cid.Cid{}, ErrBranchNotFound
	}

	return head, nil
}
