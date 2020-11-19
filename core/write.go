package core

import (
	"errors"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
)

// Write writes the contents of node to the path.
func (c *Context) Write(path string, node ipld.Node) error {
	fsnode, err := unixfs.ExtractFSNode(node)
	if err != nil {
		return err
	}

	switch fsnode.Type() {
	case unixfs.TFile:
		return c.writeFile(path, node)
	case unixfs.TDirectory:
		return c.writeDir(path, node)
	case unixfs.TSymlink:
		return c.fs.Symlink(string(fsnode.Data()), path)
	default:
		return errors.New("invalid file type")
	}
}

func (c *Context) writeFile(path string, node ipld.Node) error {
	reader, err := ufsio.NewDagReader(c.ctx, node, c.dag)
	if err != nil {
		return err
	}

	file, err := c.fs.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := reader.WriteTo(file); err != nil {
		return err
	}

	return nil
}

func (c *Context) writeDir(path string, node ipld.Node) error {
	dir, err := ufsio.NewDirectoryFromNode(c.dag, node)
	if err != nil {
		return err
	}

	if err := c.fs.MkdirAll(path, 0755); err != nil {
		return err
	}

	links, err := dir.Links(c.ctx)
	for _, link := range links {
		subnode, err := link.GetNode(c.ctx, c.dag)
		if err != nil {
			return err
		}

		subpath := c.fs.Join(path, link.Name)
		if err := c.Write(subpath, subnode); err != nil {
			return err
		}
	}

	return nil
}
