package repo

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
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
		return "", err
	}

	return Root(parent)
}

// Commit records changes in the local repo.
func Commit(ipfs *core.IpfsNode, message string) (*cid.Cid, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root, err := Root(cwd)
	if err != nil {
		return nil, err
	}

	changes, err := files.Add(ipfs, root)
	if err != nil {
		return nil, err
	}

	node, err := changes.Node()
	if err != nil {
		return nil, err
	}

	if err := ipfs.DAG.Add(context.TODO(), node); err != nil {
		return nil, err
	}

	head, err := ReadHead(root)
	if err != nil {
		return nil, err
	}

	dag, err := commit.NewCommit(message, node.Cid(), head).Node()
	if err != nil {
		return nil, err
	}

	if err := ipfs.DAG.Add(context.TODO(), dag); err != nil {
		return nil, err
	}

	id := dag.Cid()
	if err := WriteHead(root, id); err != nil {
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

	root, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	multiDir := filepath.Join(root, ".multi")
	if err := os.MkdirAll(multiDir, 0755); err != nil {
		return err
	}

	if err := WriteHead(root, id); err != nil {
		return err 
	}

	return files.Write(ipfs, root, c.Changes)
}

// Log lists change history from the current head.
func Log(ipfs *core.IpfsNode) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	head, err := ReadHead(cwd)
	if err != nil {
		return err
	}

	return commit.Log(ipfs, head)
}

// ReadHead returns the head for the repo in the given directory.
func ReadHead(dir string) (cid.Cid, error) {
	root, err := Root(dir)
	if err != nil {
		return cid.Cid{}, err
	}

	bytes, err := ioutil.ReadFile(filepath.Join(root, ".multi/HEAD"))
	if err != nil {
		return cid.Cid{}, nil
	}

	hash := string(bytes)
	return cid.Parse(hash)
}

// WriteHead sets the head for the repo in the given directory.
func WriteHead(dir string, id cid.Cid) error {
	root, err := Root(dir)
	if err != nil {
		return err
	}

	bytes := []byte(id.String())
	return ioutil.WriteFile(filepath.Join(root, ".multi/HEAD"), bytes, 0755)
}