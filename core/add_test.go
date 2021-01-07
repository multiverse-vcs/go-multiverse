package core

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/spf13/afero"
)

func TestAddFile(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "test.txt", []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(ctx, dag, "test.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	r, err := ufsio.NewDagReader(ctx, node, dag)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	if string(data) != "foo bar" {
		t.Errorf("file data does not match")
	}
}

func TestAddDir(t *testing.T) {
	fs = afero.NewMemMapFs()

	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	if err := fs.Mkdir("test", 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := filepath.Join("test", "test.txt")
	if err := afero.WriteFile(fs, path, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(ctx, dag, "test", nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	ufsdir, err := ufsio.NewDirectoryFromNode(dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = ufsdir.Find(ctx, "test.txt")
	if err != nil {
		t.Errorf("failed to find file")
	}
}

func TestFilter(t *testing.T) {
	filter := Filter{
		"*.exe",
		"baz/*",
	}

	if !filter.Match("foo.exe") {
		t.Errorf("expected filter to match")
	}

	if !filter.Match("foo/bar.exe") {
		t.Errorf("expected filter to match")
	}

	if !filter.Match("baz/bar") {
		t.Errorf("expected filter to match")
	}
}
