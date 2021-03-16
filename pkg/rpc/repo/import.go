package repo

import (
	"context"
	"errors"
	"path"

	cid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiverse-vcs/go-multiverse/internal/git"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// ImportArgs contains the args.
type ImportArgs struct {
	// Peer is the author peer ID.
	Peer string `json:"key"`
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

	priv, err := s.Keystore.Get(args.Peer)
	if err != nil {
		return err
	}

	peerID, err := peer.IDFromPrivateKey(priv)
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
	authorID, err = object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	reply.Remote = path.Join(s.Peer.Host.ID().Pretty(), args.Name)
	return s.Namesys.Publish(ctx, priv, authorID)
}
