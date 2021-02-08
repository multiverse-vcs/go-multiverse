package git

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
)

func TestImportFromURL(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	id, err := ImportFromURL(ctx, dag, "test", "https://github.com/multiverse-vcs/go-multiverse")
	if err != nil {
		t.Fatal("failed to import git repo")
	}

	repo, err := data.GetRepository(ctx, dag, id)
	if err != nil || repo == nil {
		t.Fatal("failed to get repo")
	}
}

func TestImportFromFS(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	dir, err := ioutil.TempDir("", "*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	opts := git.CloneOptions{
		URL: "https://github.com/hsoft/collapseos",
	}

	_, err = git.PlainClone(dir, false, &opts)
	if err != nil {
		t.Fatal("failed to clone repo")
	}

	id, err := ImportFromFS(ctx, dag, "test", dir)
	if err != nil {
		t.Fatalf("failed to import git repo %s", err)
	}

	repo, err := data.GetRepository(ctx, dag, id)
	if err != nil || repo == nil {
		t.Fatal("failed to get repo")
	}
}
