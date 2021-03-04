package file

import (
	"context"
	"errors"
	"strings"

	path "github.com/ipfs/go-path"
	unixfs "github.com/ipfs/go-unixfs"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// SearchArgs contains the args.
type SearchArgs struct {
	// Remote is the remote path.
	Remote string `json:"remote"`
	// Branch is the branch name.
	Branch string `json:"branch"`
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

	parts := strings.Split(args.Remote, "/")
	if len(parts) != 2 {
		return errors.New("invalid remote path")
	}

	pname := parts[0]
	rname := parts[1]

	peerID, err := peer.Decode(pname)
	if err != nil {
		return err
	}

	authorID, err := s.Namesys.Search(ctx, peerID)
	if err != nil {
		return err
	}

	author, err := object.GetAuthor(ctx, s.Peer.DAG, authorID)
	if err != nil {
		return err
	}

	repoID, ok := author.Repositories[rname]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := object.GetRepository(ctx, s.Peer.DAG, repoID)
	if err != nil {
		return err
	}

	head, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	fpath, err := path.FromSegments("/ipfs/", head.String(), "tree", args.Path)
	if err != nil {
		return err
	}

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
