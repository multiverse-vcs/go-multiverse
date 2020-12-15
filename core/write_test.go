package core

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/spf13/afero"
)

func TestWriteFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := afero.WriteFile(fs, "test.txt", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), fs, dag, "test.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := fs.Remove("test.txt"); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := Write(context.TODO(), fs, dag, "test.txt", node); err != nil {
		t.Fatalf("failed to write node")
	}

	file, err := fs.Open("test.txt")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read file")
	}

	if string(data) != "hello" {
		t.Errorf("file data does not match")
	}
}

func TestWriteDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	dag := dagutils.NewMemoryDagService()

	if err := fs.Mkdir("test", 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := filepath.Join("test", "test.txt")
	if err := afero.WriteFile(fs, path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), fs, dag, "test", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := fs.RemoveAll("test"); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := Write(context.TODO(), fs, dag, "test", node); err != nil {
		t.Fatalf("failed to write node")
	}

	if _, err := fs.Stat(path); err != nil {
		t.Fatalf("failed to lstat file")
	}
}
