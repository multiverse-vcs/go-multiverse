package core

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-ipld-multiverse"
)

// stack is a first in last out data structure.
type stack []cid.Cid

// pop removes the last item from the stack.
func (s stack) pop() (cid.Cid, stack) {
	return s[len(s)-1], s[:len(s)-1]
}

// History is a commit history iterator.
type History struct {
	core  *Core
	curr  cid.Cid
	seen  map[string]bool
	stack stack
	valid HistoryFilter
	limit HistoryFilter
}

// HistoryCallback returns an error.
type HistoryCallback func(*ipldmulti.Commit) error

// HistoryFilter returns a bool indicating if the commit should be skipped.
type HistoryFilter func(*ipldmulti.Commit) bool

// NewHistory returns a new history starting at id.
func (c *Core) NewHistory(id cid.Cid) *History {
	return &History{
		core:  c,
		seen:  map[string]bool{},
		stack: stack{id},
	}
}

// WithFilter returns a filtered history using valid and limit.
func (h *History) WithFilter(valid, limit *HistoryFilter) *History {
	h.valid = *valid
	h.limit = *limit
	return h
}

// Next returns the next commit.
func (h *History) Next(ctx context.Context) (*ipldmulti.Commit, error) {
	for {
		if len(h.stack) == 0 {
			return nil, io.EOF
		}

		h.curr, h.stack = h.stack.pop()
		if h.seen[h.curr.KeyString()] {
			continue
		}

		commit, err := h.core.Reference(ctx, path.IpfsPath(h.curr))
		if err != nil {
			return nil, err
		}

		h.seen[h.curr.KeyString()] = true
		if h.limit == nil || !h.limit(commit) {
			h.stack = append(h.stack, commit.Parents...)
		}

		if h.valid == nil || h.valid(commit) {
			return commit, nil
		}
	}
}

// ForEach walks history and invokes the call back for each commit.
func (h *History) ForEach(ctx context.Context, cb HistoryCallback) error {
	for {
		commit, err := h.Next(ctx)
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

// Flatten returns a flattened map of all commit keys in history.
func (h *History) Flatten(ctx context.Context) (map[string]bool, error) {
	for {
		_, err := h.Next(ctx)
		if err == io.EOF {
			return h.seen, nil
		}

		if err != nil {
			return nil, err
		}
	}
}
