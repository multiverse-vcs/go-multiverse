package core

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestMergeConflicts(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	base, err := Commit(context.TODO(), store, "base")
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello\nfoo\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	local, err := Commit(context.TODO(), store, "local", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := afero.WriteFile(store.Cwd, "README.md", []byte("hello\nbar\n"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	remote, err := Commit(context.TODO(), store, "remote", base)
	if err != nil {
		t.Fatalf("failed to commit")
	}

	if err := Merge(context.TODO(), store, local, remote); err != nil {
		t.Fatalf("failed to merge %s", err)
	}

	file, err := store.Cwd.Open("README.md")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read file")
	}

	expect := `hello
<<<<<<<
foo
=======
bar
>>>>>>>
`

	if string(data) != expect {
		t.Error("unexpected merge contents")
	}
}
