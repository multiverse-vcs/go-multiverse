package git

import (
	"context"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestImportFromURL(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()
	url := "https://github.com/multiverse-vcs/go-multiverse"

	_, err := ImportFromURL(ctx, mem, "test", url)
	if err != nil {
		t.Fatalf("failed to import git repo %v", err)
	}
}
