package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestMergeBase(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	base, err := Commit(context.TODO(), fs, dag, "base")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(context.TODO(), fs, dag, "local", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), fs, dag, "remote", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), dag, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != base {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseRemoteAhead(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	base, err := Commit(context.TODO(), fs, dag, "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(context.TODO(), fs, dag, "local", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), fs, dag, "remote", local)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), dag, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != local {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseLocalAhead(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	base, err := Commit(context.TODO(), fs, dag, "init")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), fs, dag, "remote", base)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	local, err := Commit(context.TODO(), fs, dag, "local", remote)
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), dag, local, remote)
	if err != nil {
		t.Fatalf("failed to get merge base")
	}

	if merge != remote {
		t.Errorf("unexpected merge base")
	}
}

func TestMergeBaseUnrelated(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	local, err := Commit(context.TODO(), fs, dag, "local")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	remote, err := Commit(context.TODO(), fs, dag, "remote")
	if err != nil {
		t.Fatalf("failed to create commit")
	}

	merge, err := MergeBase(context.TODO(), dag, local, remote)
	if merge.Defined() {
		t.Errorf("uexpected merge base")
	}
}
