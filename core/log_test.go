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

	readme := mock.fs.Join(mock.fs.Root(), "README.md")
	if err := fsutil.WriteFile(mock.fs, readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	idA, err := mock.Commit("first")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	idB, err := mock.Commit("second")
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
		t.Fatalf("failed to read log")
	}

	if !bytes.Contains(log, []byte(idA.String())) {
		t.Errorf("expected commit cid in log")
	}

	if !bytes.Contains(log, []byte(idB.String())) {
		t.Errorf("expected commit cid in log")
	}
}
