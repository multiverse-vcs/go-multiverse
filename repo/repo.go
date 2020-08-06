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
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multihash"
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
func Commit(message string) (cid.Cid, error) {
	ipfs, err := httpapi.NewLocalApi()
	if err != nil {
		return cid.Cid{}, err
	}

	node, err := Changes()
	if err != nil {
		return cid.Cid{}, err
	}

	changes, err := ipfs.Unixfs().Add(context.TODO(), node)
	if err != nil {
		return cid.Cid{}, err
	}

	commit := make(map[string]interface{})
	commit["message"] = message
	commit["changes"] = changes.Cid()

	if head, err := readHead(); err == nil {
		commit["parent"] = head
	}

	dag, err := cbornode.WrapObject(commit, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, err
	}

	if err := ipfs.Dag().Pinning().Add(context.TODO(), dag); err != nil {
		return cid.Cid{}, err
	}

	if err := writeHead(dag.Cid().Bytes()); err != nil {
		return cid.Cid{}, err
	}

	return dag.Cid(), nil
}

// Clone an existing commit into the directory.
func Clone(hash string, target string) (string, error) {
	ipfs, err := httpapi.NewLocalApi()
	if err != nil {
		return "", err
	}

	commit, err := cid.Parse(hash)
	if err != nil {
		return "", err
	}

	dag, err := ipfs.Dag().Get(context.TODO(), commit)
	if err != nil {
		return "", err
	}

	changes, _, err := dag.ResolveLink([]string{"changes"})
	if err != nil {
		return "", err
	}

	node, err := ipfs.Unixfs().Get(context.TODO(), path.IpfsPath(changes.Cid))
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

	if err := writeHead(commit.Bytes()); err != nil {
		return "", err
	}

	return dir, nil
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
		return cid.Cid{}, err
	}

	return cid.Parse(head)
}

func writeHead(cid []byte) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	defer os.Chdir(cwd)

	if err := Root(); err != nil {
		return err
	}

	return ioutil.WriteFile(".multi/HEAD", cid, 0644)
}