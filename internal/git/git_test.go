package git

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestImportFromURL(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	_, err := ImportFromURL(ctx, dag, "test", "https://github.com/multiverse-vcs/go-multiverse")
	if err != nil {
		t.Fatal("failed to import git repo")
	}
}
