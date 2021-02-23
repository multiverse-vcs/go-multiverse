package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// Fetch returns the repository at the given path.
func (s *Service) Fetch(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	params := httprouter.ParamsFromContext(ctx)
	pname := params.ByName("peer")
	rname := params.ByName("repo")

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repo)

	return nil
}
