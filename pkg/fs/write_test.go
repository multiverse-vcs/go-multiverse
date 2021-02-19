package fs

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestWriteFile(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tmp, err := ioutil.TempDir("", "unixfs-*")
	if err != nil {
		t.Fatalf("failed to create temp dir")
	}
	defer os.RemoveAll(tmp)

	original, err := ioutil.ReadFile("testdata/a.txt")
	if err != nil {
		t.Fatalf("failed to read file")
	}

	node, err := Add(ctx, dag, "testdata/a.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	path := filepath.Join(tmp, "a.txt")
	if err := Write(ctx, dag, path, node); err != nil {
		t.Fatalf("failed to write node")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file")
	}

	if !bytes.Equal(original, data) {
		t.Errorf("file data does not match")
	}
}

func TestWriteDir(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	tmp, err := ioutil.TempDir("", "unixfs-*")
	if err != nil {
		t.Fatalf("failed to create temp dir")
	}
	defer os.RemoveAll(tmp)

	node, err := Add(ctx, dag, "testdata", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	path := filepath.Join(tmp, "testdata")
	if err := Write(ctx, dag, path, node); err != nil {
		t.Fatalf("failed to write node")
	}

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatalf("failed to read dir")
	}

	if len(entries) != 6 {
		t.Fatalf("unexpected directory entries")
	}

	if entries[0].Name() != "a.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[0].IsDir() != false {
		t.Error("unexpected dir entry")
	}

	if entries[1].Name() != "b" {
		t.Error("unexpected dir entry")
	}

	if entries[1].IsDir() != true {
		t.Error("unexpected dir entry")
	}

	if entries[2].Name() != "b.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[2].IsDir() != false {
		t.Error("unexpected dir entry")
	}

	if entries[3].Name() != "l" {
		t.Error("unexpected dir entry")
	}

	if entries[3].Mode()&os.ModeSymlink == 0 {
		t.Error("unexpected dir entry")
	}

	if entries[4].Name() != "o.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[4].IsDir() != false {
		t.Error("unexpected dir entry")
	}

	if entries[5].Name() != "r.txt" {
		t.Error("unexpected dir entry")
	}

	if entries[5].IsDir() != false {
		t.Error("unexpected dir entry")
	}
}
