package core

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestWriteFile(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, "test.txt", []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), store, "test.txt", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := store.Cwd.Remove("test.txt"); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := Write(context.TODO(), store, "test.txt", node); err != nil {
		t.Fatalf("failed to write node")
	}

	file, err := store.Cwd.Open("test.txt")
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
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := store.Cwd.Mkdir("test", 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := filepath.Join("test", "test.txt")
	if err := afero.WriteFile(store.Cwd, path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := Add(context.TODO(), store, "test", nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := store.Cwd.RemoveAll("test"); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := Write(context.TODO(), store, "test", node); err != nil {
		t.Fatalf("failed to write node")
	}

	if _, err := store.Cwd.Stat(path); err != nil {
		t.Fatalf("failed to lstat file")
	}
}
