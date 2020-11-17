package object

import (
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/peer"
)

const id = "bagaybqabciqb47qmevdwdtciz43gts6ng5lbknaycovotirp76hs4t6nmoqiwbi"

var data = []byte(`{
	"date": "2020-10-25T15:26:12.168056-07:00",
	"message": "big changes",
	"parents": [{"/": "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa"}],
	"peer_id": "QmcSMDFbN4Br31GnqfhNEkqVFj3gyuVVeYriZNqY8kQpDN",
	"tree": {"/": "QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ"}
}`)

func TestDecodeCommit(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Error("failed to decode commit")
	}

	if commit.Cid() != id {
		t.Error("cid does not match")
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

	if commit.PeerID.String() != "QmcSMDFbN4Br31GnqfhNEkqVFj3gyuVVeYriZNqY8kQpDN" {
		t.Error("peer id does not match")
	}

	if commit.WorkTree.String() != "QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ" {
		t.Error("work tree does not match")
	}
}

func TestCommitResolveEmptyPath(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.Resolve([]string{})
	if err != nil {
		t.Error("failed to resolve")
	}

	other, ok := obj.(*Commit)
	if !ok {
		t.Error("expected commit")
	}

	if commit != other {
		t.Error("commit does not match")
	}

	if len(rest) != 0 {
		t.Error("rest should be empty")
	}
}

func TestCommitResolveNoSuchPath(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	_, _, err = commit.Resolve([]string{"invalidpath"})
	if err != ErrNoLink {
		t.Error("expected error no link")
	}
}

func TestCommitResolveDate(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.Resolve([]string{"date"})
	if err != nil {
		t.Error("failed to resolve date")
	}

	date, ok := obj.(time.Time)
	if !ok {
		t.Error("expected date")
	}

	if date.Format(time.RFC3339) != "2020-10-25T15:26:12-07:00" {
		t.Error("date does not match")
	}

	if len(rest) != 0 {
		t.Error("rest should be empty")
	}
}

func TestCommitResolveMessage(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.Resolve([]string{"message"})
	if err != nil {
		t.Error("failed to resolve message")
	}

	message, ok := obj.(string)
	if !ok {
		t.Error("expected message")
	}

	if message != "big changes" {
		t.Error("message does not match")
	}

	if len(rest) != 0 {
		t.Error("rest should be empty")
	}
}

func TestCommitResolvePeerID(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.Resolve([]string{"peer_id"})
	if err != nil {
		t.Error("failed to resolve message")
	}

	peerID, ok := obj.(peer.ID)
	if !ok {
		t.Error("expected peer id")
	}

	if peerID.String() != "QmcSMDFbN4Br31GnqfhNEkqVFj3gyuVVeYriZNqY8kQpDN" {
		t.Error("peer id does not match")
	}

	if len(rest) != 0 {
		t.Error("rest should be empty")
	}
}

func TestCommitResolveTree(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.Resolve([]string{"tree"})
	if err != nil {
		t.Error("failed to resolve message")
	}

	link, ok := obj.(*format.Link)
	if !ok {
		t.Error("object is not a link")
	}

	if link.Cid.String() != "QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ" {
		t.Error("link cid does not match")
	}

	if len(rest) != 0 {
		t.Error("rest should be empty")
	}
}

func TestCommitResolveParents(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.Resolve([]string{"parents"})
	if err != nil {
		t.Error("failed to resolve message")
	}

	parents, ok := obj.([]cid.Cid)
	if !ok {
		t.Error("expected parents")
	}

	if len(parents) != 1 {
		t.Error("parents does not match")
	}

	if parents[0].String() != "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa" {
		t.Error("parents does not match")
	}

	if len(rest) != 0 {
		t.Error("rest should be empty")
	}
}

func TestCommitResolveParentsIndex(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.Resolve([]string{"parents", "0"})
	if err != nil {
		t.Error("failed to resolve message")
	}

	link, ok := obj.(*format.Link)
	if !ok {
		t.Error("expected link")
	}

	if link.Cid.String() != "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa" {
		t.Error("parents does not match")
	}

	if len(rest) != 0 {
		t.Error("rest should be empty")
	}
}

func TestCommitTreeBadParams(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	tree := commit.Tree("non-empty-string", 0)
	if tree != nil {
		t.Error("expected nil tree")
	}

	tree = commit.Tree("non-empty-string", 1)
	if tree != nil {
		t.Error("expected nil tree")
	}

	tree = commit.Tree("", 0)
	if tree != nil {
		t.Error("expected nil tree")
	}
}

func TestCommitTree(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	expected := map[string]struct{}{
		"date":      {},
		"message":   {},
		"peer_id":   {},
		"tree":      {},
		"parents/0": {},
	}

	tree := commit.Tree("", -1)
	if len(tree) != len(expected) {
		t.Error("tree does not match")
	}

	for _, entry := range tree {
		if _, ok := expected[entry]; !ok {
			t.Error("tree does not match")
		}
	}
}

func TestCommitResolveLinkNoSuchLink(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	obj, rest, err := commit.ResolveLink([]string{"invalidlink"})
	if obj != nil {
		t.Error("expected obj to be nil")
	}

	if rest != nil {
		t.Error("expected rest to be nil")
	}

	if err != ErrNoLink {
		t.Error("expected err no link")
	}
}

func TestCommitResolveLinkTree(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	link, rest, err := commit.ResolveLink([]string{"tree"})
	if link == nil {
		t.Error("expected link to be non nil")
	}

	if link.Cid.String() != "QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ" {
		t.Error("link cid does not match")
	}

	if len(rest) != 0 {
		t.Error("expected rest to be empty")
	}
}

func TestCommitResolveLinkParents(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	link, rest, err := commit.ResolveLink([]string{"parents", "0"})
	if link == nil {
		t.Error("expected link to be non nil")
	}

	if link.Cid.String() != "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa" {
		t.Error("link cid does not match")
	}

	if len(rest) != 0 {
		t.Error("expected rest to be empty")
	}
}

func TestCommitLinks(t *testing.T) {
	id, err := cid.Parse(id)
	if err != nil {
		t.Fatal("failed to parse cid")
	}

	commit, err := DecodeCommit(id, data)
	if err != nil {
		t.Fatal("failed to decode commit")
	}

	links := commit.Links()
	if len(links) != 2 {
		t.Error("expected links len to be 2")
	}

	if links[0].Cid.String() != "QmQycvPQd5tAVP4Xx1dp1Yfb9tmjKQAa5uxPoTfUQr9tFZ" {
		t.Error("link cid does not match")
	}

	if links[1].Cid.String() != "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa" {
		t.Error("link cid does not match")
	}
}
