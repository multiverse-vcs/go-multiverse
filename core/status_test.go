package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestStatusRemove(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	head, err := Commit(ctx, dag, "", "init")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := fs.Remove("README.md"); err != nil {
		t.Fatalf("failed to remove readme file")
	}

	diffs, err := Status(ctx, dag, "/", head)
	if err != nil {
		t.Fatalf("failed to get status")
	}

	diff, ok := diffs["README.md"]
	if !ok {
		t.Fatalf("unexpected changes")
	}

	if diff != dagutils.Remove {
		t.Fatalf("unexpected change type")
	}
}

func TestStatusBare(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	diffs, err := Status(ctx, dag, "", cid.Cid{})
	if err != nil {
		t.Fatalf("failed to get status")
	}

	diff, ok := diffs["README.md"]
	if !ok {
		t.Fatalf("unexpected changes")
	}

	if diff != dagutils.Add {
		t.Fatalf("unexpected change type")
	}
}
