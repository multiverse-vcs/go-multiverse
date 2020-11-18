package core

import (
	"io/ioutil"
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
	"github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
)

func TestAddFile(t *testing.T) {
	mock := NewMockContext()

	path := mock.fs.Join(mock.config.Root, "test.txt")
	if err := fsutil.WriteFile(mock.fs, path, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	info, err := mock.fs.Lstat(path)
	if err != nil {
		t.Fatalf("failed to stat file")
	}

	node, err := mock.Add(path, info, nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	r, err := ufsio.NewDagReader(mock.ctx, node, mock.dag)
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
	mock := NewMockContext()

	path := mock.fs.Join(mock.config.Root, "test.txt")
	if err := fsutil.WriteFile(mock.fs, path, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	info, err := mock.fs.Lstat(mock.config.Root)
	if err != nil {
		t.Fatalf("failed to lstat")
	}

	node, err := mock.Add(mock.config.Root, info, nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	ufsdir, err := ufsio.NewDirectoryFromNode(mock.dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = ufsdir.Find(mock.ctx, path)
	if err != nil {
		t.Errorf("failed to find file")
	}
}

func TestAddSymlink(t *testing.T) {
	mock := NewMockContext()

	path := mock.fs.Join(mock.config.Root, "link")
	if err := mock.fs.Symlink("target", path); err != nil {
		t.Fatalf("failed to create symlink")
	}

	info, err := mock.fs.Lstat(path)
	if err != nil {
		t.Fatalf("failed to lstat")
	}

	node, err := mock.Add(path, info, nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	fsnode, err := unixfs.ExtractFSNode(node)
	if err != nil {
		t.Fatalf("failed to extract fsnode")
	}

	if fsnode.Type() != unixfs.TSymlink {
		t.Errorf("invalid fsnode type")
	}

	if string(fsnode.Data()) != "target" {
		t.Errorf("symlink target does not match")
	}
}
