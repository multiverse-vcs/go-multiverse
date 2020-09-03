package files

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

const BlockSize = 512

// Node represents a file or directory.
type Node struct {
	Name   string      `refmt:"name"`
	Mode   os.FileMode `refmt:"mode"`
	Files  []cid.Cid   `refmt:"files,omitempty"`
	Chunks []cid.Cid   `refmt:"chunks,omitempty"`
}

func init() {
	cbornode.RegisterCborType(Node{})
}

// Get returns the file node with the given CID.
func Get(ipfs *core.IpfsNode, id cid.Cid) (*Node, error) {
	dag, err := ipfs.DAG.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	var node Node
	if err := cbornode.DecodeInto(dag.RawData(), &node); err != nil {
		return nil, err
	}

	return &node, nil
}

// Add creates a file or directory node from the given path.
func Add(ipfs *core.IpfsNode, path string) (*Node, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return AddDirectory(ipfs, path, stat)
	}

	return AddFile(ipfs, path, stat)
}

// AddFile creates a file node from the given path.
func AddFile(ipfs *core.IpfsNode, path string, info os.FileInfo) (*Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	splitter := NewSizeSplitter(file, BlockSize)
	size := info.Size() / BlockSize
	if (info.Size() % BlockSize) > 0 {
		size = size + 1
	}

	chunks := make([]cid.Cid, size)
	for i := range chunks {
		chunk, err := splitter.NextBytes()
		if err != nil {
			return nil, err
		}

		hash, err := multihash.Sum(chunk, multihash.SHA2_256, -1)
		if err != nil {
			return nil, err
		}

		block, err := blocks.NewBlockWithCid(chunk, cid.NewCidV1(cid.Raw, hash))
		if  err != nil {
			return nil, err
		}

		if err := ipfs.Blocks.AddBlock(block); err != nil {
			return nil, err
		}

		chunks[i] = block.Cid()
	}

	return &Node{Name: info.Name(), Mode: info.Mode(), Chunks: chunks}, nil
}

// AddDirectory creates a directory node from the given path.
func AddDirectory(ipfs *core.IpfsNode, path string, info os.FileInfo) (*Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	files := make([]cid.Cid, len(names))
	for i, name := range names {
		node, err := Add(ipfs, filepath.Join(path, name))
		if err != nil {
			return nil, err
		}

		dag, err := node.Node()
		if err != nil {
			return nil, err
		}

		if err := ipfs.DAG.Add(context.TODO(), dag); err != nil {
			return nil, err
		}

		files[i] = dag.Cid()
	}

	return &Node{Name: info.Name(), Mode: info.Mode(), Files: files}, nil
}

// Write copies the file or directory to the local file system.
func Write(ipfs *core.IpfsNode, path string, id cid.Cid) error {
	node, err := Get(ipfs, id)
	if err != nil {
		return err
	}

	if node.Mode.IsDir() {
		return WriteDirectory(ipfs, path, node)
	}

	return WriteFile(ipfs, path, node)
}

// WriteFile copies the file to the local file system.
func WriteFile(ipfs *core.IpfsNode, path string, node *Node) error {
	file, err := os.Create(filepath.Join(path, node.Name))
	if err != nil {
		return err
	}

	defer file.Close()

	for i, chunk := range node.Chunks {
		block, err := ipfs.Blocks.GetBlock(context.TODO(), chunk)
		if err != nil {
			return err
		}

		offset := (int64) (i * BlockSize)
		if _, err := file.WriteAt(block.RawData(), offset); err != nil {
			return err
		}
	}

	return nil
}

// WriteDirectory copies the directory to the local file system.
func WriteDirectory(ipfs *core.IpfsNode, path string, node *Node) error {
	path = filepath.Join(path, node.Name)
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	for _, file := range node.Files {
		if err := Write(ipfs, path, file); err != nil {
			return err
		}
	}

	return nil
}

// Node returns an ipld representation of the node.
func (n *Node) Node() (format.Node, error) {
	return cbornode.WrapObject(n, multihash.SHA2_256, -1)
}