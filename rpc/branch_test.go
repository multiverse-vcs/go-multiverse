package rpc

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/peer"
)

func TestListBranches(t *testing.T) {
	ctx := context.Background()
	node := peer.NewMock(t, ctx)

	json, err := ioutil.ReadFile("testdata/repository.json")
	if err != nil {
		t.Fatal("failed to read json")
	}

	repo, err := data.RepositoryFromJSON(json)
	if err != nil {
		t.Fatal("failed to parse repository json")
	}

	id, err := data.AddRepository(ctx, node.Dag(), repo)
	if err != nil {
		t.Fatal("failed to create repository")
	}
	node.Config().Author.Repositories["test"] = id

	client, err := connect(node)
	if err != nil {
		t.Fatal("failed to connect to rpc server")
	}

	args := BranchArgs{
		Name: "test",
	}

	var reply BranchReply
	if err := client.Call("Service.ListBranches", &args, &reply); err != nil {
		t.Fatal("failed to call rpc")
	}

	if len(reply.Branches) != 1 {
		t.Error("unexpected branches")
	}

	if reply.Branches["default"] != repo.Branches["default"] {
		t.Error("unexpected branches")
	}
}
