package git

import (
	"context"
	"testing"
	"path/filepath"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestImportFromURL(t *testing.T) {
	ctx := context.Background()
	mem := dagutils.NewMemoryDagService()

	path, err := filepath.Abs("./../../")
	if err != nil {
		t.Fatal("failed to get absolute path")
	}

	_, err = ImportFromFS(ctx, mem, "test", path)
	if err != nil {
		t.Fatal("failed to import git repo")
	}
}
