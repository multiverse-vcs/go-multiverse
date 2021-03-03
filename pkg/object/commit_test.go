package object

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ipfs/go-merkledag/dagutils"
)

func TestCommitRoundtrip(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	data, err := os.ReadFile("testdata/commit.json")
	if err != nil {
		t.Fatal("failed to read file")
	}

	commit, err := CommitFromJSON(data)
	if err != nil {
		t.Fatal("failed to decode commit json")
	}

	id, err := AddCommit(ctx, dag, commit)
	if err != nil {
		t.Fatal("failed to add commit to dag")
	}

	commit, err = GetCommit(ctx, dag, id)
	if err != nil {
		t.Fatal("failed to add commit to dag")
	}

	if commit.Message != "big changes" {
		t.Error("message does not match")
	}

	if commit.Date.Format(time.RFC3339) != "2020-10-25T15:26:12-07:00" {
		t.Error("date does not match")
	}

	if len(commit.Parents) != 1 {
		t.Error("parents does not match")
	}

	if commit.Parents[0].String() != "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa" {
		t.Error("parents does not match")
	}

	if commit.Tree.String() != "QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ" {
		t.Error("work tree does not match")
	}

	meta, ok := commit.Metadata["foo"]
	if !ok || meta != "bar" {
		t.Error("metadata does not match")
	}
}

func TestCommitParentLinks(t *testing.T) {
	data, err := os.ReadFile("testdata/commit.json")
	if err != nil {
		t.Fatal("failed to read file")
	}

	commit, err := CommitFromJSON(data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	links := commit.ParentLinks()
	if len(links) != 1 {
		t.Error("expected parent links len to be 1")
	}

	if links[0].Cid.String() != "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa" {
		t.Error("parent link cid does not match")
	}
}
