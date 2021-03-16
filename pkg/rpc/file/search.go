package file

import (
	"context"

	path "github.com/ipfs/go-path"
	unixfs "github.com/ipfs/go-unixfs"

	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
)

// SearchArgs contains the args.
type SearchArgs struct {
	// Path is the file path.
	Path string `json:"path"`
}

// SearchReply contains the reply.
type SearchReply struct {
	// Content contains file content.
	Content string `json:"content"`
	// Entries contains directory entries.
	Entries []*fs.DirEntry `json:"entries"`
	// IsDir specifies if the file is a directory.
	IsDir bool `json:"is_dir"`
}

// Search returns the contents of a file at the given remote path.
func (s *Service) Search(args *SearchArgs, reply *SearchReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fpath := path.FromString(args.Path)
	fnode, err := s.Resolver.ResolvePath(ctx, fpath)
	if err != nil {
		return err
	}

	fsnode, err := unixfs.ExtractFSNode(fnode)
	if err != nil {
		return err
	}

	switch {
	case fsnode.IsDir():
		tree, err := fs.Ls(ctx, s.Peer.DAG, fnode.Cid())
		if err != nil {
			return err
		}

		reply.Entries = tree
	default:
		blob, err := fs.Cat(ctx, s.Peer.DAG, fnode.Cid())
		if err != nil {
			return err
		}

		reply.Content = blob
	}

	reply.IsDir = fsnode.IsDir()
	return nil
}
