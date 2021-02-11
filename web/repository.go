package web

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/ipfs/go-path"
	"github.com/julienschmidt/httprouter"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

type Repository Server

var readmePattern = regexp.MustCompile(`(?i)^readme.*`)

func (s *Repository) Index(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()
	dag := s.node.Dag()

	params := httprouter.ParamsFromContext(ctx)
	peerID := params.ByName("peer_id")
	name := params.ByName("name")
	refs := params.ByName("refs")
	head := params.ByName("head")
	file := params.ByName("file")

	pid, err := peer.Decode(peerID)
	if err != nil {
		return nil, err
	}

	author, err := s.node.Authors().Search(ctx, pid)
	if err != nil {
		return nil, err
	}

	repoID, ok := author.Repositories[name]
	if !ok {
		return nil, errors.New("repository does not exist")
	}

	repo, err := data.GetRepository(ctx, dag, repoID)
	if err != nil {
		return nil, err
	}

	fpath, err := path.FromSegments("/ipfs/", repoID.String(), refs, head, "tree", file)
	if err != nil {
		return nil, err
	}

	fnode, err := s.node.ResolvePath(ctx, fpath)
	if err != nil {
		return nil, err
	}

	dir, err := unixfs.IsDir(ctx, dag, fnode.Cid())
	if err != nil {
		return nil, err
	}

	var blob string
	var tree []*unixfs.DirEntry

	var readmeBlob string
	var readmeName string

	switch {
	case dir:
		tree, err = unixfs.Ls(ctx, dag, fnode.Cid())
		if err != nil {
			return nil, err
		}

		entry, err := unixfs.Find(ctx, dag, fnode.Cid(), readmePattern)
		if err != nil || entry == nil {
			break
		}

		readmeName = entry.Name
		readmeBlob, err = unixfs.Cat(ctx, dag, entry.Cid)
		if err != nil {
			return nil, err
		}
	default:
		blob, err = unixfs.Cat(ctx, dag, fnode.Cid())
		if err != nil {
			return nil, err
		}
	}

	return &ViewModel{
		Name: "repository.html",
		Data: map[string]interface{}{
			"RepoID":     repoID,
			"PeerID":     peerID,
			"Name":       name,
			"Repo":       repo,
			"Refs":       refs,
			"File":       file,
			"Head":       head,
			"Blob":       blob,
			"Tree":       tree,
			"ReadmeBlob": readmeBlob,
			"ReadmeName": readmeName,
		},
	}, nil
}
