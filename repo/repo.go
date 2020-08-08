package repo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/multiverse/commit"
)

// Walk parent directories until repo root is found.
func Root() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	multi := filepath.Join(cwd, ".multi")
	if dirExists(multi) {
		return nil
	}

	parent, err := filepath.Abs("..")
	if err != nil {
		return err
	}

	if parent == cwd {
		return fmt.Errorf("Repo not found")
	}

	if err := os.Chdir(parent); err != nil {
		return err
	}

	return Root()
}

// Record changes in the local repo.
func Commit(message string) (*commit.Commit, error) {
	ipfs, err := httpapi.NewLocalApi()
	if err != nil {
		return nil, err
	}

	node, err := Changes()
	if err != nil {
		return nil, err
	}

	changes, err := ipfs.Unixfs().Add(context.TODO(), node)
	if err != nil {
		return nil, err
	}

	head, err := readHead()
	if err != nil {
		return nil, err
	}

	c := commit.NewCommit(message, changes.Cid(), head)

	dag, err := c.Node(ipfs)
	if err != nil {
		return nil, err
	}

	if err := ipfs.Dag().Pinning().Add(context.TODO(), dag); err != nil {
		return nil, err
	}

	if err := writeHead(dag.Cid()); err != nil {
		return nil, err
	}

	c.Id = dag.Cid()
	return c, nil
}

// Clone an existing commit into the directory.
func Clone(id cid.Cid, target string) (string, error) {
	ipfs, err := httpapi.NewLocalApi()
	if err != nil {
		return "", err
	}

	c, err := commit.Get(ipfs, id)
	if err != nil {
		return "", err
	}

	node, err := ipfs.Unixfs().Get(context.TODO(), path.IpfsPath(c.Changes))
	if err != nil {
		return "", err
	}

	dir, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}

	if err := files.WriteTo(node, dir); err != nil {
		return "", err
	}

	if err := os.Chdir(dir); err != nil {
		return "", err
	}

	if err := os.MkdirAll(".multi", os.ModePerm); err != nil {
		return "", err
	}

	if err := writeHead(id); err != nil {
		return "", err
	}

	return dir, nil
}

// List change history from the current head.
func Log() error {
	ipfs, err := httpapi.NewLocalApi()
	if err != nil {
		return nil
	}

	head, err := readHead()
	if err != nil {
		return err
	}

	c, err := commit.Get(ipfs, head)
	if err != nil {
		return err
	}

	for {
		fmt.Println(c.String())
		if !c.Parent.Defined() {
			return nil
		}

		c, err = commit.Get(ipfs, c.Parent)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns a node containing the local changes.
func Changes() (files.Node, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	defer os.Chdir(cwd)

	if err := Root(); err != nil {
		return nil, err
	}

	root, err := os.Open(".")
	if err != nil {
		return nil, err
	}

	defer root.Close()

	stat, err := root.Stat()
	if err != nil {
		return nil, err
	}

	filter, err := files.NewFilter("", []string{".multi"}, true)
	if err != nil {
		return nil, err
	}

	return files.NewSerialFileWithFilter(".", filter, stat)
}

func dirExists(path string) bool {
	dir, err := os.Open(path)
	if err != nil {
		return false
	}

	defer dir.Close()

	stat, err := dir.Stat()
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func readHead() (cid.Cid, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return cid.Cid{}, err
	}

	defer os.Chdir(cwd)

	if err := Root(); err != nil {
		return cid.Cid{}, err
	}

	head, err := ioutil.ReadFile(".multi/HEAD")
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

	defer os.Chdir(cwd)

	if err := Root(); err != nil {
		return err
	}

	return ioutil.WriteFile(".multi/HEAD", id.Bytes(), 0644)
}