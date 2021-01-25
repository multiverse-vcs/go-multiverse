package rpc

import (
	"context"
	"testing"

	"github.com/ipfs/go-datastore"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/peer"
)

var repositoryJSON = []byte(`{
	"name": "test",
	"branches": {
		"default": {"/": "bagaybqabciqeutn2u7n3zuk5b4ykgfwpkekb7ctgnlwik5zfr6bcukvknj2jtpa"}
	},
	"tags": {}
}`)

func TestListBranches(t *testing.T) {
	ctx := context.Background()

	dstore := datastore.NewMapDatastore()
	store := data.NewStore(dstore)

	mock, err := peer.Mock(ctx, dstore)
	if err != nil {
		t.Fatalf("failed to create peer")
	}

	repo, err := data.RepositoryFromJSON(repositoryJSON)
	if err != nil {
		t.Fatalf("failed to parse repository json")
	}

	id, err := data.AddRepository(ctx, mock, repo)
	if err != nil {
		t.Fatalf("failed to create repository %s", err)
	}

	if err := store.PutCid(repo.Name, id); err != nil {
		t.Fatalf("failed to put cid in store")
	}

	client, err := connect(mock, store)
	if err != nil {
		t.Fatalf("failed to connect to rpc server")
	}

	args := BranchArgs{
		Name: repo.Name,
	}

	var reply BranchReply
	if err := client.Call("Service.ListBranches", &args, &reply); err != nil {
		t.Fatalf("failed to call rpc: %s", err)
	}

	if len(reply.Branches) != 1 {
		t.Errorf("unexpected branches")
	}

	if reply.Branches["default"] != repo.Branches["default"] {
		t.Errorf("unexpected branches")
	}
}
