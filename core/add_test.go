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
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "test.txt", []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), fs, dag, "test.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	r, err := ufsio.NewDagReader(context.TODO(), node, dag)
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
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := fs.Mkdir("test", 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := filepath.Join("test", "test.txt")
	if err := afero.WriteFile(fs, path, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), fs, dag, "test", nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	ufsdir, err := ufsio.NewDirectoryFromNode(dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = ufsdir.Find(context.TODO(), "test.txt")
	if err != nil {
		t.Errorf("failed to find file")
	}
}
