package core

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/yondero/go-ipld-multiverse"
)

// History is a commit history iterator.
type History struct {
	core  *Core
	seen  map[string]bool
	stack []cid.Cid
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
		stack: []cid.Cid{id},
	}
}

// NewFilterHistory returns a new filtered history starting at id and stopping at limit.
func (c *Core) NewFilterHistory(id cid.Cid, valid *HistoryFilter, limit *HistoryFilter) *History {
	return &History{
		core:  c,
		seen:  map[string]bool{},
		stack: []cid.Cid{id},
		valid: *valid,
		limit: *limit,
	}
}

// Next returns the next commit.
func (h *History) Next(ctx context.Context) (*ipldmulti.Commit, error) {
	for {
		index := len(h.stack) - 1
		if index < 0 {
			return nil, io.EOF
		}

		id := h.stack[index]

		h.stack = h.stack[:index]
		if h.seen[id.KeyString()] {
			continue
		}

		node, err := h.core.Api.Dag().Get(ctx, id)
		if err != nil {
			return nil, err
		}

		commit, ok := node.(*ipldmulti.Commit)
		if !ok {
			return nil, ErrInvalidRef
		}

		h.seen[id.KeyString()] = true
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