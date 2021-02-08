package rpc

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/multiverse-vcs/go-multiverse/data"
)

func TestLog(t *testing.T) {
	ctx := context.Background()

	node, err := makeNode(ctx)
	if err != nil {
		t.Fatal("failed to create peer node")
	}

	json, err := ioutil.ReadFile("testdata/repository.json")
	if err != nil {
		t.Fatal("failed to read json")
	}

	repo, err := data.RepositoryFromJSON(json)
	if err != nil {
		t.Fatal("failed to parse repository json")
	}

	id, err := data.AddRepository(ctx, node, repo)
	if err != nil {
		t.Fatal("failed to add repository")
	}
	node.Config().Author.Repositories["test"] = id

	json, err = ioutil.ReadFile("testdata/commit.json")
	if err != nil {
		t.Fatal("failed to read json")
	}

	commit, err := data.CommitFromJSON(json)
	if err != nil {
		t.Fatal("failed to parse commit json")
	}

	head, err := data.AddCommit(ctx, node, commit)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	client, err := makeClient(node)
	if err != nil {
		t.Fatal("failed to connect to rpc server")
	}

	args := LogArgs{
		Name:   "test",
		Branch: "default",
		Limit:  1,
	}

	var reply LogReply
	if err := client.Call("Service.Log", &args, &reply); err != nil {
		t.Fatal("failed to call rpc")
	}

	if len(reply.IDs) != 1 {
		t.Fatal("unexpected ids")
	}

	if reply.IDs[0] != head {
		t.Error("unexpected log id")
	}

	if len(reply.Commits) != 1 {
		t.Fatal("unexpected commits")
	}
}
