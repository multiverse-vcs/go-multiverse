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

	path := mock.Fs.Join(mock.Fs.Root(), "test.txt")
	if err := fsutil.WriteFile(mock.Fs, path, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := mock.Add(path, nil)
	if err != nil {
		t.Fatalf("failed to add file")
	}

	r, err := ufsio.NewDagReader(mock, node, mock.Dag)
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

	dir := mock.Fs.Join(mock.Fs.Root(), "test")
	if err := mock.Fs.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to mkdir")
	}

	path := mock.Fs.Join(dir, "test.txt")
	if err := fsutil.WriteFile(mock.Fs, path, []byte("foo bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	node, err := mock.Add(dir, nil)
	if err != nil {
		t.Fatalf("failed to add")
	}

	ufsdir, err := ufsio.NewDirectoryFromNode(mock.Dag, node)
	if err != nil {
		t.Fatalf("failed to read node")
	}

	_, err = ufsdir.Find(mock, "test.txt")
	if err != nil {
		t.Errorf("failed to find file")
	}
}

func TestAddSymlink(t *testing.T) {
	mock := NewMockContext()

	path := mock.Fs.Join(mock.Fs.Root(), "link")
	if err := mock.Fs.Symlink("target", path); err != nil {
		t.Fatalf("failed to create symlink")
	}

	node, err := mock.Add(path, nil)
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
