package core

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
)

func TestLog(t *testing.T) {
	mock := NewMockContext()

	readme := mock.fs.Join(mock.config.Root, "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	commit1, err := mock.Commit("first")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	commit2, err := mock.Commit("second")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	r, w := io.Pipe()
	go func() {
		mock.Log(w)
		w.Close()
	}()

	log, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read log: %s", err)
	}

	cid1 := commit1.Cid().String()
	if !bytes.Contains(log, []byte(cid1)) {
		t.Errorf("expected commit cid in log")
	}

	cid2 := commit2.Cid().String()
	if !bytes.Contains(log, []byte(cid2)) {
		t.Errorf("expected commit cid in log")
	}
}
