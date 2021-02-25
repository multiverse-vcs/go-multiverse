package http

import (
	"encoding/json"
	"errors"
	"net/http"

	path "github.com/ipfs/go-path"
	unixfs "github.com/ipfs/go-unixfs"
	"github.com/julienschmidt/httprouter"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/fs"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// File returns file contents from the repository tree.
func (s *Service) File(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	highlight := req.URL.Query().Get("highlight")
	params := httprouter.ParamsFromContext(ctx)
	pname := params.ByName("peer")
	rname := params.ByName("repo")
	bname := params.ByName("branch")
	fname := params.ByName("file")

	peerID, err := peer.Decode(pname)
	if err != nil {
		return err
	}

	authorID, err := s.Namesys.Resolve(ctx, peerID)
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

	head, ok := repo.Branches[bname]
	if !ok {
		return errors.New("branch does not exist")
	}

	fpath, err := path.FromSegments("/ipfs/", head.String(), "tree", fname)
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tree)
	case highlight != "":
		blob, err := fs.Cat(ctx, s.Peer.DAG, fnode.Cid())
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		Highlight(fname, blob, highlight, w)
	default:
		blob, err := fs.Cat(ctx, s.Peer.DAG, fnode.Cid())
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(blob))
	}

	return nil
}
