package fs

import (
	"context"
	"sort"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	unixfs "github.com/ipfs/go-unixfs"
	io "github.com/ipfs/go-unixfs/io"
)

// DirEntry contains info about a file.
type DirEntry struct {
	// Cid is the id of the entry.
	Cid cid.Cid `json:"cid"`
	// Name is the name of the file.
	Name string `json:"name"`
	// IsDir indicates if the file is a directory.
	IsDir bool `json:"is_dir"`
	// Size is the size in bytes of the file.
	Size uint64 `json:"size"`
}

// ByTypeAndName is used to sort dir entries.
type ByTypeAndName []*DirEntry

// Len returns the length of the slice.
func (a ByTypeAndName) Len() int {
	return len(a)
}

// Swap swaps the indices of two entries.
func (a ByTypeAndName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less returns true if entry at index i is less than entry at index j.
func (a ByTypeAndName) Less(i, j int) bool {
	if a[i].IsDir != a[j].IsDir {
		return a[i].IsDir
	}

	return a[i].Name < a[j].Name
}

// Ls returns the contents of the directory with the given CID.
func Ls(ctx context.Context, dag ipld.DAGService, id cid.Cid) ([]*DirEntry, error) {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	dir, err := io.NewDirectoryFromNode(dag, node)
	if err != nil {
		return nil, err
	}

	links, err := dir.Links(ctx)
	if err != nil {
		return nil, err
	}

	var list []*DirEntry
	for _, l := range links {
		info, err := Stat(ctx, dag, l)
		if err != nil {
			return nil, err
		}

		list = append(list, info)
	}

	sort.Sort(ByTypeAndName(list))
	return list, nil
}

// Stat returns file info for the given link.
func Stat(ctx context.Context, dag ipld.DAGService, link *ipld.Link) (*DirEntry, error) {
	node, err := dag.Get(ctx, link.Cid)
	if err != nil {
		return nil, err
	}

	fsnode, err := unixfs.ExtractFSNode(node)
	if err != nil {
		return nil, err
	}

	return &DirEntry{
		Cid:   link.Cid,
		Name:  link.Name,
		IsDir: fsnode.IsDir(),
		Size:  fsnode.FileSize(),
	}, nil
}
