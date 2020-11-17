package core

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-chunker"
	"github.com/ipfs/go-ipfs-files"
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

// Add adds a file to the merkle dag.
func (c *Context) Add(file files.Node) (ipld.Node, error) {
	adder, err := c.newAdder()
	if err != nil {
		return nil, err
	}

	node, err := adder.addNode(file)
	if err != nil {
		return nil, err
	}

	if err := c.dag.Add(c.ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

type adder struct {
	ctx     context.Context
	dag     ipld.DAGService
	builder cid.Builder
}

func (c *Context) newAdder() (*adder, error) {
	prefix, err := merkledag.PrefixForCidVersion(1)
	if err != nil {
		return nil, err
	}

	prefix.MhType = multihash.SHA2_256
	prefix.MhLength = -1

	return &adder{
		ctx:     c.ctx,
		dag:     c.dag,
		builder: &prefix,
	}, nil
}

func (adder *adder) addNode(file files.Node) (ipld.Node, error) {
	defer file.Close()

	switch node := file.(type) {
	case files.Directory:
		dir := ufsio.NewDirectory(adder.dag)
		if err := adder.addEntries(node.Entries(), dir); err != nil {
			return nil, err
		}

		return dir.GetNode()
	case *files.Symlink:
		data, err := unixfs.SymlinkData(node.Target)
		if err != nil {
			return nil, err
		}

		return merkledag.NodeWithData(data), nil
	case files.File:
		return adder.addFile(node)
	default:
		return nil, errors.New("invalid file type")
	}
}

func (adder *adder) addFile(file files.File) (ipld.Node, error) {
	params := helpers.DagBuilderParams{
		Dagserv:    adder.dag,
		CidBuilder: adder.builder,
		Maxlinks:   helpers.DefaultLinksPerBlock,
	}

	chunker, err := chunk.FromString(file, DefaultChunker)
	if err != nil {
		return nil, err
	}

	helper, err := params.New(chunker)
	if err != nil {
		return nil, err
	}

	return balanced.Layout(helper)
}

func (adder *adder) addEntries(entries files.DirIterator, dir ufsio.Directory) error {
	if !entries.Next() {
		return entries.Err()
	}

	node, err := adder.addNode(entries.Node())
	if err != nil {
		return err
	}

	if err := adder.dag.Add(adder.ctx, node); err != nil {
		return err
	}

	if err := dir.AddChild(adder.ctx, entries.Name(), node); err != nil {
		return err
	}

	return adder.addEntries(entries, dir)
}
