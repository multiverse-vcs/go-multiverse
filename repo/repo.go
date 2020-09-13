package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/yondero/multiverse/commit"
	"github.com/yondero/multiverse/file"
	"github.com/yondero/multiverse/ipfs"
)

// Repo contains repository info.
type Repo struct {
	path string
}

// Open returns a repo if one exists in the current or parent directories.
func Open(path string) (*Repo, error) {
	stat, err := os.Stat(filepath.Join(path, ".multi"))
	if err == nil && stat.IsDir() {
		return &Repo{path}, nil
	}

	parent, err := filepath.Abs(filepath.Join(path, ".."))
	if err != nil {
		return nil, err
	}

	if parent == path {
		return nil, err
	}

	return Open(parent)
}

// Clone creates a repo from an existing commit.
func Clone(ipfs *ipfs.Node, id cid.Cid, root string) (*Repo, error) {
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

	if err := file.Write(ipfs, c.Changes, root); err != nil {
		return nil, err
	}

	return &Repo{path: root}, nil
}

// Commit records changes in the local repo.
func (r *Repo) Commit(ipfs *ipfs.Node, message string) (*commit.Commit, error) {
	f, err := file.NewFile(r.path)
	if err != nil {
		return nil, err
	}

	if err := f.Add(ipfs); err != nil {
		return nil, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	c := commit.NewCommit(message, f.ID, head)
	if err := c.Add(ipfs); err != nil {
		return nil, err
	}

	path := filepath.Join(r.path, ".multi", "HEAD")
	if err := ioutil.WriteFile(path, []byte(c.ID.String()), 0755); err != nil {
		return nil, err 
	}

	return c, nil
}

// Log prints the commit history from the given CID.
func (r *Repo) Log(ipfs *ipfs.Node, id cid.Cid) error {
	c, err := commit.Get(ipfs, id)
	if err != nil {
		return err
	}

	fmt.Println(c.String())
	if c.Parent.Defined() {
		return r.Log(ipfs, c.Parent)
	}

	return nil
}

// Head returns the parent CID of the repo.
func (r *Repo) Head() (cid.Cid, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(r.path, ".multi", "HEAD"))
	if err != nil {
		return cid.Cid{}, nil
	}

	return cid.Parse(string(bytes))
}

// String returns a string representation of the repo.
func (r *Repo) String() string {
	return fmt.Sprintf("repo %s", r.path)
}