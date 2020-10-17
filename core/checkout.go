package core

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/yondero/go-ipld-multiverse"
)

// Checkout copies the tree of the commit with the given path to the local repo directory.
func (c *Core) Checkout(ctx context.Context, ref path.Path) error {
	p, err := c.Api.ResolvePath(ctx, ref)
	if err != nil {
		return err
	}

	if p.Cid().Type() != ipldmulti.CommitCodec {
		return ErrInvalidRef
	}

	node, err := c.Api.Unixfs().Get(ctx, path.Join(p, "tree"))
	if err != nil {
		return err
	}

	if err := checkoutNode(node, c.Config.Path); err != nil {
		return err
	}

	c.Config.Head = p.Root()
	return c.Config.Write()
}

func checkoutNode(node files.Node, root string) error {
	dir, ok := node.(files.Directory)
	if ok {
		return checkoutDirectory(dir, root)
	}

	file, ok := node.(files.File)
	if ok {
		return checkoutFile(file, root)
	}

	return ErrInvalidFile
}

func checkoutFile(node files.File, root string) error {
	b, err := ioutil.ReadAll(node)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(root, b, 0644)
}

func checkoutDirectory(node files.Directory, root string) error {
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	return checkoutEntries(node.Entries(), root)
}

func checkoutEntries(entries files.DirIterator, root string) error {
	if !entries.Next() {
		return entries.Err()
	}

	path := filepath.Join(root, entries.Name())
	if err := checkoutNode(entries.Node(), path); err != nil {
		return err
	}

	return checkoutEntries(entries, root)
}