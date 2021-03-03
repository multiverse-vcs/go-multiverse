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
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

func TestTree(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()

	treeO, err := fs.Add(ctx, mem, "testdata/o", nil)
	if err != nil {
		t.Fatal("failed to add dir")
	}

	commitO := object.NewCommit()
	commitO.Tree = treeO.Cid()

	o, err := object.AddCommit(ctx, mem, commitO)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	treeA, err := fs.Add(ctx, mem, "testdata/a", nil)
	if err != nil {
		t.Fatal("failed to add dir")
	}

	commitA := object.NewCommit()
	commitA.Tree = treeA.Cid()

	a, err := object.AddCommit(ctx, mem, commitA)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	treeB, err := fs.Add(ctx, mem, "testdata/b", nil)
	if err != nil {
		t.Fatal("failed to add dir")
	}

	commitB := object.NewCommit()
	commitB.Tree = treeB.Cid()

	b, err := object.AddCommit(ctx, mem, commitB)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	merge, err := Tree(ctx, mem, o, a, b)
	if err != nil {
		t.Fatalf("failed to merge %s", err)
	}

	ufsdir, err := ufsio.NewDirectoryFromNode(mem, merge)
	if err != nil {
		t.Fatal("failed to read node")
	}

	file, err := ufsdir.Find(ctx, "list.txt")
	if err != nil {
		t.Error("failed to find file")
	}

	r, err := ufsio.NewDagReader(ctx, file, mem)
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
