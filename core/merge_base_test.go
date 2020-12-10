package core

import (
	"context"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestMergeBase(t *testing.T) {
	store, err := storage.NewStore(afero.NewMemMapFs(), "/")
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	base, err := Commit(context.TODO(), store, "base")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(context.TODO(), store, "local", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), store, "remote", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), store, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != base {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseRemoteAhead(t *testing.T) {
	store, err := storage.NewStore(afero.NewMemMapFs(), "/")
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	base, err := Commit(context.TODO(), store, "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(context.TODO(), store, "local", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), store, "remote", local)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), store, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != local {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseLocalAhead(t *testing.T) {
	store, err := storage.NewStore(afero.NewMemMapFs(), "/")
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	base, err := Commit(context.TODO(), store, "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), store, "remote", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(context.TODO(), store, "local", remote)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), store, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != remote {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseUnrelated(t *testing.T) {
	store, err := storage.NewStore(afero.NewMemMapFs(), "/")
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	local, err := Commit(context.TODO(), store, "local")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), store, "remote")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), store, local, remote)
	if merge.Defined() {
		t.Errorf("uexpected merge base")
	}
}
