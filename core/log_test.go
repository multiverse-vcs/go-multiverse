package core

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestLog(t *testing.T) {
	c, err := NewMockContext()
	if err != nil {
		t.Fatalf("failed to create context")
	}

	commit1, err := c.Commit("first")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	commit2, err := c.Commit("second")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	r, w := io.Pipe()
	go func() {
		c.Log(w)
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
