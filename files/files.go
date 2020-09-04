package files

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/go-ipld-format"
	"github.com/multiformats/go-multihash"
)

// BlockSize is the maximum size of file chunks.
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
func Add(ipfs *core.IpfsNode, path string) (format.Node, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return addDirectory(ipfs, path, stat)
	}

	return addFile(ipfs, path, stat)
}

// Write copies the file or directory to the local file system.
func Write(ipfs *core.IpfsNode, path string, id cid.Cid) error {
	node, err := Get(ipfs, id)
	if err != nil {
		return err
	}

	if node.Mode.IsDir() {
		return writeDirectory(ipfs, path, node)
	}

	return writeFile(ipfs, path, node)
}

func addFile(ipfs *core.IpfsNode, path string, info os.FileInfo) (format.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// calculate number of chunks to reduce allocations
	size := info.Size() / BlockSize
	if (info.Size() % BlockSize) > 0 {
		size = size + 1
	}

	splitter := NewSizeSplitter(file, BlockSize)
	chunks := make([]cid.Cid, 0, size)
	for range chunks {
		block, err := splitter.NextBlock()
		if err != nil {
			return nil, err
		}

		if err := ipfs.Blocks.AddBlock(block); err != nil {
			return nil, err
		}

		chunks = append(chunks, block.Cid())
	}

	node := &Node{Name: info.Name(), Mode: info.Mode(), Chunks: chunks}
	return addNode(ipfs, node)
}

func addDirectory(ipfs *core.IpfsNode, path string, info os.FileInfo) (format.Node, error) {
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

	files := make([]cid.Cid, 0, len(names))
	for _, name := range names {
		node, err := Add(ipfs, filepath.Join(path, name))
		if err != nil {
			return nil, err
		}

		files = append(files, node.Cid())
	}

	node := &Node{Name: info.Name(), Mode: info.Mode(), Files: files}
	return addNode(ipfs, node)
}

func addNode(ipfs *core.IpfsNode, node *Node) (format.Node, error) {
	dag, err := cbornode.WrapObject(node, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	if err := ipfs.DAG.Add(context.TODO(), dag); err != nil {
		return nil, err
	}

	return dag, nil
}

func writeFile(ipfs *core.IpfsNode, path string, node *Node) error {
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

func writeDirectory(ipfs *core.IpfsNode, path string, node *Node) error {
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
