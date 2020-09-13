package file

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-block-format"
	"github.com/multiformats/go-multihash"
	"github.com/yondero/multiverse/ipfs"
)

// BlockSize is the maximum size of blocks.
const BlockSize = 512

// File contains file metadata and contents.
type File struct {
	ID     cid.Cid        `json:"-"`
	Name   string         `json:"name"`
	Mode   os.FileMode    `json:"mode"`
	Links  []cid.Cid      `json:"links"`
	Blocks []blocks.Block `json:"-"`
}

// NewFile creates a new file from the given path.
func NewFile(path string) (*File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.Mode().IsDir() {
		return newDirectoryFile(path, info)
	}

	return newRegularFile(path, info)
}

// Get returns the file node with the given CID.
func Get(ipfs *ipfs.Node, id cid.Cid) (*File, error) {
	b, err := ipfs.Blocks.GetBlock(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	f := File{ID: id}
	if err := json.Unmarshal(b.RawData(), &f); err != nil {
		return nil, err
	}

	return &f, nil
}

// Write copies the file or directory to the local file system.
func Write(ipfs *ipfs.Node, id cid.Cid, path string) error {
	f, err := Get(ipfs, id)
	if err != nil {
		return err
	}

	if f.Mode.IsDir() {
		return writeDirectoryFile(ipfs, path, f)
	}

	return writeRegularFile(ipfs, path, f)
}

// Add persists the File to the blockstore.
func (f *File) Add(ipfs *ipfs.Node) error {
	b, err := f.Block()
	if err != nil {
		return err
	}

	f.ID = b.Cid()
	if err := ipfs.Blocks.AddBlock(b); err != nil {
		return err
	}

	return nil
}

// Block returns a block representation of the File.
func (f *File) Block() (blocks.Block, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	return blocks.NewBlockWithCid(data, cid.NewCidV1(cid.Raw, hash))
}

func newRegularFile(path string, info os.FileInfo) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	size := info.Size() / BlockSize
	if (info.Size() % BlockSize) > 0 {
		size = size + 1
	}

	f := File{Name: info.Name(), Mode: info.Mode()}
	f.Blocks = make([]blocks.Block, 0, size)
	f.Links = make([]cid.Cid, 0, size)

	for i := int64(0); i < size; i++ {
		data := make([]byte, BlockSize)

		s, err := io.ReadFull(file, data)
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, err
		}

		hash, err := multihash.Sum(data[:s], multihash.SHA2_256, -1)
		if err != nil {
			return nil, err
		}

		block, err := blocks.NewBlockWithCid(data[:s], cid.NewCidV1(cid.Raw, hash))
		if err != nil {
			return nil, err
		}

		f.Blocks = append(f.Blocks, block)
		f.Links = append(f.Links, block.Cid())
	}

	return &f, nil
}

func newDirectoryFile(path string, info os.FileInfo) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	f := File{Name: info.Name(), Mode: info.Mode()}
	f.Blocks = make([]blocks.Block, 0, len(names))
	f.Links = make([]cid.Cid, 0, len(names))

	for _, name := range names {
		child, err := NewFile(filepath.Join(path, name))
		if err != nil {
			return nil, err
		}

		block, err := child.Block()
		if err != nil {
			return nil, err
		}

		f.Blocks = append(f.Blocks, block)
		f.Links = append(f.Links, block.Cid())
	}

	return &f, nil
}

func writeRegularFile(ipfs *ipfs.Node, path string, f *File) error {
	file, err := os.Create(filepath.Join(path, f.Name))
	if err != nil {
		return err
	}

	defer file.Close()

	for i, id := range f.Links {
		block, err := ipfs.Blocks.GetBlock(context.TODO(), id)
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

func writeDirectoryFile(ipfs *ipfs.Node, path string, f *File) error {
	path = filepath.Join(path, f.Name)
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	for _, id := range f.Links {
		if err := Write(ipfs, id, path); err != nil {
			return err
		}
	}

	return nil
}