package core

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestWalk(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "README.md", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	idA, err := Commit(ctx, dag, "", nil, "first")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	idB, err := Commit(ctx, dag, "", nil, "second", idA)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	history, err := Walk(ctx, dag, idB, nil)
	if err != nil {
		t.Fatalf("failed to walk")
	}

	if len(history) != 2 {
		t.Fatalf("cids do not match")
	}

	if _, ok := history[idA.KeyString()]; !ok {
		t.Errorf("cids do not match")
	}

	if _, ok := history[idB.KeyString()]; !ok {
		t.Errorf("cids do not match")
	}
}
