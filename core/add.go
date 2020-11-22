package core

import (
	"errors"
	"os"

	"github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/importer/balanced"
	"github.com/ipfs/go-unixfs/importer/helpers"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/sabhiram/go-gitignore"
)

// DefaultChunker is the name of the default chunker algorithm.
const DefaultChunker = "buzhash"

// Add creates a node from the file at path and adds it to the merkle dag.
func (c *Context) Add(path string, filter *ignore.GitIgnore) (ipld.Node, error) {
	stat, err := c.Fs.Lstat(path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsRegular():
		return c.addFile(path)
	case mode&os.ModeSymlink != 0:
		return c.addSymlink(path)
	case mode.IsDir():
		return c.addDir(path, filter)
	default:
		return nil, errors.New("invalid file type")
	}
}

func (c *Context) addFile(path string) (ipld.Node, error) {
	params := helpers.DagBuilderParams{
		Dagserv:    c.Dag,
		CidBuilder: merkledag.V1CidPrefix(),
		Maxlinks:   helpers.DefaultLinksPerBlock,
	}

	file, err := c.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	chunker, err := chunk.FromString(file, DefaultChunker)
	if err != nil {
		return nil, err
	}

	helper, err := params.New(chunker)
	if err != nil {
		return nil, err
	}

	node, err := balanced.Layout(helper)
	if err != nil {
		return nil, err
	}

	return node, c.Dag.Add(c, node)
}

func (c *Context) addSymlink(path string) (ipld.Node, error) {
	target, err := c.Fs.Readlink(path)
	if err != nil {
		return nil, err
	}

	data, err := unixfs.SymlinkData(target)
	if err != nil {
		return nil, err
	}

	node := merkledag.NodeWithData(data)
	return node, c.Dag.Add(c, node)
}

func (c *Context) addDir(path string, filter *ignore.GitIgnore) (ipld.Node, error) {
	entries, err := c.Fs.ReadDir(path)
	if err != nil {
		return nil, err
	}

	dir := ufsio.NewDirectory(c.Dag)
	for _, info := range entries {
		subpath := c.Fs.Join(path, info.Name())
		if filter != nil && filter.MatchesPath(subpath) {
			continue
		}

		subnode, err := c.Add(subpath, filter)
		if err != nil {
			return nil, err
		}

		if err := dir.AddChild(c, info.Name(), subnode); err != nil {
			return nil, err
		}
	}

	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}

	return node, c.Dag.Add(c, node)
}
