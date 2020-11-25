package core

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestAddFile(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, "test.txt", []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), store, "test.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	r, err := ufsio.NewDagReader(context.TODO(), node, store.Dag)
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
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := store.Cwd.Mkdir("test", 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := filepath.Join("test", "test.txt")
	if err := afero.WriteFile(store.Cwd, path, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), store, "test", nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	ufsdir, err := ufsio.NewDirectoryFromNode(store.Dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = ufsdir.Find(context.TODO(), "test.txt")
	if err != nil {
		t.Errorf("failed to find file")
	}
}
