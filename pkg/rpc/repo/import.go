package repo

import (
	"context"
	"errors"
	"path"

	cid "github.com/ipfs/go-cid"

	"github.com/multiverse-vcs/go-multiverse/internal/git"
	"github.com/multiverse-vcs/go-multiverse/internal/p2p"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// ImportArgs contains the args.
type ImportArgs struct {
	// Name is the repository name.
	Name string `json:"name"`
	// URL is the repository url.
	URL string `json:"url"`
	// Path is the repository directory.
	Path string `json:"path"`
}

// ImportReply contains the reply
type ImportReply struct {
	// Remote is the repository path
	Remote string `json:"remote"`
}

// Import imports an external repository.
func (s *Service) Import(args *ImportArgs, reply *ImportReply) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if args.Name == "" {
		return errors.New("name cannot be empty")
	}

	key, err := p2p.DecodeKey(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	author := s.Config.Author
	if _, ok := author.Repositories[args.Name]; ok {
		return errors.New("repository already exists")
	}

	var repoID cid.Cid
	switch {
	case args.URL != "":
		repoID, err = git.ImportFromURL(ctx, s.Peer.DAG, args.Name, args.URL)
		if err != nil {
			return err
		}
	case args.Path != "":
		repoID, err = git.ImportFromFS(ctx, s.Peer.DAG, args.Name, args.Path)
		if err != nil {
			return err
		}
	default:
		return errors.New("import path or url must be set")
	}

	author.Repositories[args.Name] = repoID
	if err := s.Config.Write(); err != nil {
		return err
	}

	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	reply.Remote = path.Join(s.Peer.Host.ID().Pretty(), args.Name)
	return s.Namesys.Publish(ctx, key, authorID)
}
