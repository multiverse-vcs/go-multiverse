package repo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipld-format"
	"github.com/yondero/multiverse/commit"
	"github.com/yondero/multiverse/files"
)

// Root walks parent directories until the root is found.
func Root(path string) (string, error) {
	stat, err := os.Stat(filepath.Join(path, ".multi"))
	if err == nil && stat.IsDir() {
		return path, nil
	}

	parent, err := filepath.Abs(filepath.Join(path, ".."))
	if err != nil {
		return "", err
	}

	if parent == path {
		return "", fmt.Errorf("Repo not found")
	}

	return Root(parent)
}

// Commit records changes in the local repo.
func Commit(ipfs *core.IpfsNode, message string) (*cid.Cid, error) {
	changes, err := Changes(ipfs)
	if err != nil {
		return nil, err
	}

	if err := ipfs.DAG.Add(context.TODO(), changes); err != nil {
		return nil, err
	}

	head, err := readHead()
	if err != nil {
		return nil, err
	}

	c := commit.NewCommit(message, changes.Cid(), head)

	dag, err := c.Node()
	if err != nil {
		return nil, err
	}

	if err := ipfs.DAG.Add(context.TODO(), dag); err != nil {
		return nil, err
	}

	id := dag.Cid()
	if err := writeHead(id); err != nil {
		return nil, err
	}

	return &id, nil
}

// Clone copies an existing commit into the directory.
func Clone(ipfs *core.IpfsNode, id cid.Cid, path string) error {
	c, err := commit.Get(ipfs, id)
	if err != nil {
		return err
	}

	dir, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0644); err != nil {
		return err
	}

	if err := files.Write(ipfs, dir, c.Changes); err != nil {
		return err
	}

	return nil
}

// Log lists change history from the current head.
func Log(ipfs *core.IpfsNode) error {
	head, err := readHead()
	if err != nil {
		return err
	}

	for head.Defined() {		
		c, err := commit.Get(ipfs, head)
		if err != nil {
			return err
		}

		head = c.Parent
		fmt.Println(c.String())
	}

	return nil
}

// Changes creates a node containing the local changes.
func Changes(ipfs *core.IpfsNode) (format.Node, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root, err := Root(cwd)
	if err != nil {
		return nil, err
	}

	node, err := files.Add(ipfs, root)
	if err != nil {
		return nil, err
	}

	return node.Node()
}

func readHead() (cid.Cid, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return cid.Cid{}, err
	}

	root, err := Root(cwd)
	if err != nil {
		return cid.Cid{}, err
	}

	head, err := ioutil.ReadFile(filepath.Join(root, ".multi/HEAD"))
	if err != nil {
		return cid.Cid{}, nil
	}

	return cid.Parse(head)
}

func writeHead(id cid.Cid) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	root, err := Root(cwd)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(root, ".multi/HEAD"), id.Bytes(), 0644)
}
