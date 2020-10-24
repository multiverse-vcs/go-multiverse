package core

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/yondero/go-ipld-multiverse"
)

// CommitIter is a commit interator.
type CommitIter struct {
	core  *Core
	seen  map[string]bool
	stack []cid.Cid
	valid CommitFilter
	limit CommitFilter
}

// CommitFilter is used to filter commits in the iterator.
type CommitFilter func(*ipldmulti.Commit) bool

// NewCommitIter returns a new commit iterator.
func (c *Core) NewCommitIter(id cid.Cid) *CommitIter {
	return &CommitIter{
		core:  c,
		seen:  map[string]bool{},
		stack: []cid.Cid{id},
	}
}

// WithFilter returns a commit iterator with a filter function.
func (i *CommitIter) WithFilter(valid *CommitFilter) *CommitIter {
	i.valid = *valid
	return i
}

// WithLimit returns a commit iterator with a limit function.
func (i *CommitIter) WithLimit(limit *CommitFilter) *CommitIter {
	i.limit = *limit
	return i
}

// Next returns the next commit in the repo history.
func (i *CommitIter) Next(ctx context.Context) (*ipldmulti.Commit, error) {
	for {
		index := len(i.stack) - 1
		if index < 0 {
			return nil, io.EOF
		}

		id := i.stack[index]

		i.stack = i.stack[:index]
		if i.seen[id.KeyString()] {
			continue
		}

		node, err := i.core.Api.Dag().Get(ctx, id)
		if err != nil {
			return nil, err
		}

		commit, ok := node.(*ipldmulti.Commit)
		if !ok {
			return nil, ErrInvalidRef
		}

		i.seen[id.KeyString()] = true
		if i.limit != nil && i.limit(commit) {
			continue
		}

		for _, p := range commit.Parents {
			i.stack = append(i.stack, p)
		}

		if i.valid == nil || i.valid(commit) {
			return commit, nil
		}
	}
}

// ForEach walks history and invokes the call back for each commit.
func (i *CommitIter) ForEach(ctx context.Context, cb func(*ipldmulti.Commit) error) error {
	for {
		commit, err := i.Next(ctx)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if err = cb(commit); err != nil {
			return err
		}
	}
}
