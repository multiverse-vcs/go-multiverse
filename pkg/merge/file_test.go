package merge

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
)

func TestFile(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()

	nodeO, err := fs.Add(ctx, mem, "testdata/o/list.txt", nil)
	if err != nil {
		t.Fatal("failed to add file")
	}

	nodeA, err := fs.Add(ctx, mem, "testdata/a/list.txt", nil)
	if err != nil {
		t.Fatal("failed to add file")
	}

	nodeB, err := fs.Add(ctx, mem, "testdata/b/list.txt", nil)
	if err != nil {
		t.Fatal("failed to add file")
	}

	merge, err := File(ctx, mem, nodeO.Cid(), nodeA.Cid(), nodeB.Cid())
	if err != nil {
		t.Fatal("failed to merge")
	}

	r, err := ufsio.NewDagReader(ctx, merge, mem)
	if err != nil {
		t.Fatal("failed to read node")
	}

	result, err := io.ReadAll(r)
	if err != nil {
		t.Fatal("failed to read node")
	}

	expect, err := os.ReadFile("testdata/merge.txt")
	if err != nil {
		t.Fatal("failed to read file")
	}

	// fix for carriage returns on windows
	result = bytes.ReplaceAll(result, []byte("\r"), nil)
	expect = bytes.ReplaceAll(expect, []byte("\r"), nil)

	if !bytes.Equal(result, expect) {
		t.Error("unexpected merge result")
	}
}
