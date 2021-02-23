package http

import (
	"errors"
	"io"
	"net/http"

	car "github.com/ipld/go-car"
	"github.com/julienschmidt/httprouter"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/multiverse-vcs/go-multiverse/pkg/p2p"
)

// Push updates a repository branch.
func (s *Service) Push(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	params := httprouter.ParamsFromContext(ctx)
	pname := params.ByName("peer")
	rname := params.ByName("repo")
	bname := params.ByName("branch")

	key, err := p2p.DecodeKey(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	peerID, err := peer.Decode(pname)
	if err != nil {
		return err
	}

	if !peerID.MatchesPrivateKey(key) {
		return errors.New("private key does not match")
	}

	author := s.Config.Author
	repoID, ok := author.Repositories[rname]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := object.GetRepository(ctx, s.Peer.DAG, repoID)
	if err != nil {
		return err
	}

	r, err := car.NewCarReader(req.Body)
	if err != nil {
		return err
	}

	if len(r.Header.Roots) != 1 {
		return errors.New("unexpected header roots")
	}

	// load blocks slowly or badger will return an error
	for {
		block, err := r.Next()
		if err == io.EOF {
			break
		}
		
		if err != nil {
			return err
		}
		
		if err := s.Peer.Blocks.Put(block); err != nil {
			return err
		}
	}

	// TODO use merge base to check if new root is valid

	repo.Branches[bname] = r.Header.Roots[0]
	if repo.DefaultBranch == "" {
		repo.DefaultBranch = bname
	}

	repoID, err = object.AddRepository(ctx, s.Peer.DAG, repo)
	if err != nil {
		return err
	}

	author.Repositories[rname] = repoID
	if err := s.Config.Write(); err != nil {
		return err
	}

	authorID, err := object.AddAuthor(ctx, s.Peer.DAG, author)
	if err != nil {
		return err
	}

	if err := s.Namesys.Publish(ctx, key, authorID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
