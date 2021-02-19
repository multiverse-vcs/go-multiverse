package fs

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestLs(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	node, err := Add(ctx, dag, "testdata", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	entries, err := Ls(ctx, dag, node.Cid())
	if err != nil {
		t.Fatalf("failed to ls")
	}

	if len(entries) != 6 {
		t.Fatalf("unexpected entries")
	}

	if entries[0].Name != "b" {
		t.Error("unexpected dir entry")
	}

	if entries[0].IsDir != true {
		t.Error("unexpected dir entry")
	}

	if entries[1].Name != "a.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[1].IsDir != false {
		t.Error("unexpected dir entry")
	}

	if entries[2].Name != "b.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[2].IsDir != false {
		t.Error("unexpected dir entry")
	}

	if entries[3].Name != "l" {
		t.Error("unexpected dir entry")
	}

	if entries[3].IsDir != false {
		t.Error("unexpected dir entry")
	}

	if entries[4].Name != "o.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[4].IsDir != false {
		t.Error("unexpected dir entry")
	}

	if entries[5].Name != "r.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[5].IsDir != false {
		t.Error("unexpected dir entry")
	}
}
