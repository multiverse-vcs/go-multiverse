package core

import (
	"io/ioutil"
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
)

func TestWriteFile(t *testing.T) {
	mock := NewMockContext()

	path := mock.Fs.Join(mock.Fs.Root(), "test.txt")
	if err := fsutil.WriteFile(mock.Fs, path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := mock.Add(path, nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := mock.Fs.Remove(path); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := mock.Write(path, node); err != nil {
		t.Fatalf("failed to write node")
	}

	file, err := mock.Fs.Open(path)
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

	path := mock.Fs.Join(mock.Fs.Root(), "link")
	if err := mock.Fs.Symlink("target", path); err != nil {
		t.Fatalf("failed to create symlink")
	}

	node, err := mock.Add(path, nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	if err := mock.Fs.Remove(path); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := mock.Write(path, node); err != nil {
		t.Fatalf("failed to write node")
	}

	target, err := mock.Fs.Readlink(path)
	if err != nil {
		t.Fatalf("failed to read symlink")
	}

	if target != "target" {
		t.Errorf("symlink target does not match")
	}
}

func TestWriteDir(t *testing.T) {
	mock := NewMockContext()

	dir := mock.Fs.Join(mock.Fs.Root(), "test")
	if err := mock.Fs.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := mock.Fs.Join(dir, "test.txt")
	if err := fsutil.WriteFile(mock.Fs, path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := mock.Add(dir, nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	if err := fsutil.RemoveAll(mock.Fs, dir); err != nil {
		t.Fatalf("failed to remove file")
	}

	if err := mock.Write(dir, node); err != nil {
		t.Fatalf("failed to write node")
	}

	if _, err := mock.Fs.Lstat(path); err != nil {
		t.Fatalf("failed to lstat file")
	}
}
