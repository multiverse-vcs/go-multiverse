package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipld-format"
	"github.com/yondero/multiverse/commit"
	"github.com/yondero/multiverse/files"
)

// Repo contains repository info.
type Repo struct {
	Root string
}

// Open returns a repo if one exists in the current or parent directories.
func Open(root string) (*Repo, error) {
	stat, err := os.Stat(filepath.Join(root, ".multi"))
	if err == nil && stat.IsDir() {
		return &Repo{Root: root}, nil
	}

	parent, err := filepath.Abs(filepath.Join(root, ".."))
	if err != nil {
		return nil, err
	}

	if parent == root {
		return nil, err
	}

	return Open(parent)
}

// Clone creates a repo from an existing commit.
func Clone(ipfs *core.IpfsNode, id cid.Cid, root string) (*Repo, error) {
	c, err := commit.Get(ipfs, id)
	if err != nil {
		return nil, err
	}

	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(root, ".multi")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "HEAD")
	if err := ioutil.WriteFile(path, []byte(id.String()), 0755); err != nil {
		return nil, err 
	}

	if err := files.Write(ipfs, root, c.Changes); err != nil {
		return nil, err
	}

	return &Repo{Root: root}, nil
}

// Commit records changes in the local repo.
func (r *Repo) Commit(ipfs *core.IpfsNode, message string) (format.Node, error) {
	node, err := files.Add(ipfs, r.Root)
	if err != nil {
		return nil, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	dag, err := commit.Add(ipfs, message, node.Cid(), head)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(r.Root, ".multi/HEAD")
	if err := ioutil.WriteFile(path, []byte(dag.Cid().String()), 0755); err != nil {
		return nil, err 
	}

	return dag, nil
}

// Log lists change history from the current head.
func (r *Repo) Log(ipfs *core.IpfsNode) error {
	head, err := r.Head()
	if err != nil {
		return err
	}

	return commit.Log(ipfs, head)
}

// Head returns the CID of the current head of the repo.
func (r *Repo) Head() (cid.Cid, error) {
	path := filepath.Join(r.Root, ".multi/HEAD")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return cid.Cid{}, nil
	}

	hash := string(bytes)
	return cid.Parse(hash)
}
