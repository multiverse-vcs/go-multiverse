package fs

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestCat(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	original, err := ioutil.ReadFile("testdata/a.txt")
	if err != nil {
		t.Fatalf("failed to read file")
	}

	node, err := Add(ctx, dag, "testdata/a.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	text, err := Cat(ctx, dag, node.Cid())
	if err != nil {
		t.Fatal("failed to cat file")
	}

	if string(original) != text {
		t.Error("unexpected file contents")
	}
}
