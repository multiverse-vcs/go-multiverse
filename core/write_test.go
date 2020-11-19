package core

import (
	"io/ioutil"
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
)

func TestWriteFile(t *testing.T) {
	mock := NewMockContext()

	path := mock.fs.Join(mock.config.Root, "test.txt")
	if err := fsutil.WriteFile(mock.fs, path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := mock.Add(path, nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := mock.fs.Remove(path); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := mock.Write(path, node); err != nil {
		t.Fatalf("failed to write node")
	}

	file, err := mock.fs.Open(path)
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

func TestWriteSymlink(t *testing.T) {
	mock := NewMockContext()

	path := mock.fs.Join(mock.config.Root, "link")
	if err := mock.fs.Symlink("target", path); err != nil {
		t.Fatalf("failed to create symlink")
	}

	node, err := mock.Add(path, nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	if err := mock.fs.Remove(path); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := mock.Write(path, node); err != nil {
		t.Fatalf("failed to write node")
	}

	target, err := mock.fs.Readlink(path)
	if err != nil {
		t.Fatalf("failed to read symlink")
	}

	if target != "target" {
		t.Errorf("symlink target does not match")
	}
}

func TestWriteDir(t *testing.T) {
	mock := NewMockContext()

	dir := mock.fs.Join(mock.config.Root, "test")
	if err := mock.fs.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := mock.fs.Join(dir, "test.txt")
	if err := fsutil.WriteFile(mock.fs, path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := mock.Add(dir, nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := fsutil.RemoveAll(mock.fs, dir); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := mock.Write(dir, node); err != nil {
		t.Fatalf("failed to write node")
	}

	if _, err := mock.fs.Lstat(path); err != nil {
		t.Fatalf("failed to lstat file")
	}
}
