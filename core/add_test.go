package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
	ufsio "github.com/ipfs/go-unixfs/io"
)

func TestAdd(t *testing.T) {
	mock := NewMockContext()

	path1 := mock.fs.Join(mock.config.Root, "1")
	if err := fsutil.WriteFile(mock.fs, path1, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	path2 := mock.fs.Join(mock.config.Root, "2")
	if err := mock.fs.Symlink(path1, path2); err != nil {
		t.Fatalf("failed to write file")
	}

	path3 := mock.fs.Join(mock.config.Root, "3")
	if err := mock.fs.MkdirAll(path3, 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	info, err := mock.fs.Lstat(mock.config.Root)
	if err != nil {
		t.Fatalf("failed to lstat")
	}

	adder, err := mock.NewAdder()
	if err != nil {
		t.Fatalf("failed to create adder")
	}

	node, err := adder.Add(mock.config.Root, info)
	if err != nil {
		t.Fatalf("failed to add file: %s", err)
	}

	dir, err := ufsio.NewDirectoryFromNode(mock.dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = dir.Find(mock.ctx, path1)
	if err != nil {
		t.Errorf("failed to find file")
	}

	_, err = dir.Find(mock.ctx, path2)
	if err != nil {
		t.Errorf("failed to find file")
	}

	_, err = dir.Find(mock.ctx, path3)
	if err != nil {
		t.Errorf("failed to find file")
	}
}
