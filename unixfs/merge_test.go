package unixfs

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	ufsio "github.com/ipfs/go-unixfs/io"
)

func TestMerge(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	result, err := ioutil.ReadFile("testdata/r.txt")
	if err != nil {
		t.Fatalf("failed to read file")
	}

	nodeO, err := Add(ctx, dag, "testdata/o.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	nodeA, err := Add(ctx, dag, "testdata/a.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	nodeB, err := Add(ctx, dag, "testdata/b.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	merge, err := Merge(ctx, dag, nodeO.Cid(), nodeA.Cid(), nodeB.Cid())
	if err != nil {
		t.Fatalf("failed to merge")
	}

	r, err := ufsio.NewDagReader(ctx, merge, dag)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	if !bytes.Equal(result, data) {
		t.Errorf("unexpected merge result")
	}
}
