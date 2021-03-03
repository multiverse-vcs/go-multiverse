package object

import (
	"context"
	"os"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestAuthorRoundtrip(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	data, err := os.ReadFile("testdata/author.json")
	if err != nil {
		t.Fatal("failed to read file")
	}

	author, err := AuthorFromJSON(data)
	if err != nil {
		t.Fatal("failed to decode author json")
	}

	id, err := AddAuthor(ctx, dag, author)
	if err != nil {
		t.Fatal("failed to add author to dag")
	}

	author, err = GetAuthor(ctx, dag, id)
	if err != nil {
		t.Fatal("failed to get author from dag")
	}

	if len(author.Metadata) != 1 {
		t.Error("unexpected metadata")
	}

	meta, ok := author.Metadata["foo"]
	if !ok || meta != "bar" {
		t.Error("unexpected metadata value")
	}

	if len(author.Repositories) != 1 {
		t.Error("unexpected repositories")
	}

	repo, ok := author.Repositories["test"]
	if !ok || repo.String() != "bafyreib2rnmsqouz67uvb4jcjsqdvsmakdn3zrpswt4ud7aegbfohyrkbe" {
		t.Error("unexpected repository value")
	}

	if len(author.Following) != 1 {
		t.Fatal("unexpected following")
	}

	if author.Following[0].String() != "12D3KooWGacxGyqrDFTkCW9Br1TmesJ9DB84Hch5Mz9uZSbK9BeQ" {
		t.Error("unexpected following value")
	}
}
