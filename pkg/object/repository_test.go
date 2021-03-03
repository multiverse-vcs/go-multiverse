package object

import (
	"context"
	"os"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestRepositoryRoundtrip(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	data, err := os.ReadFile("testdata/repository.json")
	if err != nil {
		t.Fatal("failed to read file")
	}

	repo, err := RepositoryFromJSON(data)
	if err != nil {
		t.Fatal("failed to decode repo json")
	}

	id, err := AddRepository(ctx, dag, repo)
	if err != nil {
		t.Fatal("failed to add repo to dag")
	}

	repo, err = GetRepository(ctx, dag, id)
	if err != nil {
		t.Fatal("failed to get repo from dag")
	}

	if repo.DefaultBranch != "default" {
		t.Error("default branch does not match")
	}

	if len(repo.Branches) != 1 {
		t.Fatal("unexpected branches")
	}

	branch, ok := repo.Branches["default"]
	if !ok || branch.String() != "bafyreieo2mhnqyqntenwyndzxoovw5nhbpit727kjrl3mjbyb5nv6zs2pu" {
		t.Error("unexpected branch value")
	}

	if len(repo.Tags) != 1 {
		t.Fatal("unexpected tags")
	}

	tag, ok := repo.Tags["v0.0.1"]
	if !ok || tag.String() != "bafyreieo2mhnqyqntenwyndzxoovw5nhbpit727kjrl3mjbyb5nv6zs3pu" {
		t.Error("unexpected tag value")
	}

	if len(repo.Metadata) != 1 {
		t.Fatal("unexpected metadata")
	}

	meta, ok := repo.Metadata["foo"]
	if !ok || meta != "bar" {
		t.Error("unexpected metadata value")
	}
}
