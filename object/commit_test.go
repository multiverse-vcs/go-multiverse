package object

import (
	"testing"
	"time"

	"github.com/ipfs/go-ipld-cbor"
)

var data = []byte(`{
	"date": "2020-10-25T15:26:12.168056-07:00",
	"message": "big changes",
	"parents": [{"/": "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa"}],
	"peer_id": "QmcSMDFbN4Br31GnqfhNEkqVFj3gyuVVeYriZNqY8kQpDN",
	"tree": {"/": "QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ"}
}`)

func TestCommitFromJSON(t *testing.T) {
	commit, err := CommitFromJSON(data)
	if err != nil {
		t.Fatalf("failed to decode commit json")
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
}

func TestCommitFromCBOR(t *testing.T) {
	commit, err := CommitFromJSON(data)
	if err != nil {
		t.Fatalf("failed to decode commit")
	}

	data, err := cbornode.DumpObject(commit)
	if err != nil {
		t.Fatalf("failed to encode commit")
	}

	commit, err = CommitFromCBOR(data)
	if err != nil {
		t.Fatalf("failed to decode commit")
	}

	if commit.Message != "big changes" {
		t.Error("message does not match")
	}

	if commit.Date.Format(time.RFC3339) != "2020-10-25T15:26:12-07:00" {
		t.Errorf("date does not match")
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
}

func TestCommitParentLinks(t *testing.T) {
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
