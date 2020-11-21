package core

import (
	"testing"

	"github.com/ipfs/go-cid"
)

func TestMergeBase(t *testing.T) {
	mock := NewMockContext()

	if err := mock.fs.MkdirAll(mock.fs.Root(), 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	base, err := mock.Commit("base")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := mock.Commit("local")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	mock.cfg.Head = base
	mock.cfg.Base = base

	remote, err := mock.Commit("remote")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := mock.MergeBase(local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != base {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseRemoteAhead(t *testing.T) {
	mock := NewMockContext()

	if err := mock.fs.MkdirAll(mock.fs.Root(), 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	if _, err := mock.Commit("init"); err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := mock.Commit("local")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := mock.Commit("remote")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := mock.MergeBase(local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != local {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseLocalAhead(t *testing.T) {
	mock := NewMockContext()

	if err := mock.fs.MkdirAll(mock.fs.Root(), 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	if _, err := mock.Commit("init"); err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := mock.Commit("remote")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := mock.Commit("local")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := mock.MergeBase(local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != remote {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseUnrelated(t *testing.T) {
	mock := NewMockContext()

	if err := mock.fs.MkdirAll(mock.fs.Root(), 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	local, err := mock.Commit("local")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	mock.cfg.Head = cid.Cid{}
	mock.cfg.Base = cid.Cid{}

	remote, err := mock.Commit("remote")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := mock.MergeBase(local, remote)
	if merge.Defined() {
		t.Errorf("uexpected merge base")
	}
}
