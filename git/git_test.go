package git

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/multiverse-vcs/go-multiverse/data"
)

func TestImportFromURL(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	id, err := ImportFromURL(ctx, dag, "go-multiverse", "https://github.com/multiverse-vcs/go-multiverse")
	if err != nil {
		t.Fatal("failed to import git repo")
	}

	repo, err := data.GetRepository(ctx, dag, id)
	if err != nil {
		t.Fatal("failed to get repo")
	}

	if repo.Name != "go-multiverse" {
		t.Error("unexpected repo name")
	}
}