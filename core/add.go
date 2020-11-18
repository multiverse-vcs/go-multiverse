package core

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/importer/balanced"
	"github.com/ipfs/go-unixfs/importer/helpers"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/multiformats/go-multihash"
)

// DefaultChunker is the name of the default chunker algorithm.
const DefaultChunker = "buzhash"

// Adder is used to add files to the merkle dag.
type Adder struct {
	ctx    context.Context
	dag    ipld.DAGService
	fs     billy.Filesystem
	prefix *cid.Prefix
}

// NewAdder returns an adder with default settings.
func (c *Context) NewAdder() (*Adder, error) {
	prefix, err := merkledag.PrefixForCidVersion(1)
	if err != nil {
		return nil, err
	}

	prefix.MhType = multihash.SHA2_256
	prefix.MhLength = -1

	return &Adder{
		ctx:    c.ctx,
		dag:    c.dag,
		fs:     c.fs,
		prefix: &prefix,
	}, nil
}

// Add creates a node from the file at path and adds it to the merkle dag.
func (adder *Adder) Add(path string, stat os.FileInfo) (ipld.Node, error) {
	switch mode := stat.Mode(); {
	case mode.IsRegular():
		return adder.addFile(path)
	case mode.IsDir():
		return adder.addDir(path)
	case mode&os.ModeSymlink != 0:
		return adder.addSymlink(path)
	default:
		return nil, errors.New("invalid file type")
	}
}

func (adder *Adder) add(r io.Reader) (ipld.Node, error) {
	params := helpers.DagBuilderParams{
		Dagserv:    adder.dag,
		CidBuilder: adder.prefix,
		Maxlinks:   helpers.DefaultLinksPerBlock,
	}

	chunker, err := chunk.FromString(r, DefaultChunker)
	if err != nil {
		return nil, err
	}

	helper, err := params.New(chunker)
	if err != nil {
		return nil, err
	}

	return balanced.Layout(helper)
}

func (adder *Adder) addFile(path string) (ipld.Node, error) {
	file, err := adder.fs.Open(path)
	if err != nil {
		return nil, err
	}

	node, err := adder.add(file)
	if err != nil {
		return nil, err
	}

	return node, adder.dag.Add(adder.ctx, node)
}

func (adder *Adder) addSymlink(path string) (ipld.Node, error) {
	target, err := adder.fs.Readlink(path)
	if err != nil {
		return nil, err
	}

	data, err := unixfs.SymlinkData(target)
	if err != nil {
		return nil, err
	}

	node := merkledag.NodeWithData(data)
	return node, adder.dag.Add(adder.ctx, node)
}

func (adder *Adder) addDir(path string) (ipld.Node, error) {
	entries, err := adder.fs.ReadDir(path)
	if err != nil {
		return nil, err
	}

	dir := ufsio.NewDirectory(adder.dag)
	for _, info := range entries {
		subpath := adder.fs.Join(path, info.Name())
		subnode, err := adder.Add(subpath, info)
		if err != nil {
			return nil, err
		}

		if err := dir.AddChild(adder.ctx, info.Name(), subnode); err != nil {
			return nil, err
		}
	}

	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}

	return node, adder.dag.Add(adder.ctx, node)
}
