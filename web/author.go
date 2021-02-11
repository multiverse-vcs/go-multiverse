package web

import (
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/multiverse-vcs/go-multiverse/data"
)

type Author Server

func (s *Author) Index(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()
	cfg := s.node.Config()
	dag := s.node.Dag()

	params := httprouter.ParamsFromContext(ctx)
	peerID := params.ByName("peer_id")
	selfID := s.node.ID().Pretty()

	if peerID == "" {
		peerID = selfID
	}

	pid, err := peer.Decode(peerID)
	if err != nil {
		return nil, err
	}

	author, err := s.node.Authors().Search(ctx, pid)
	if err != nil {
		return nil, err
	}

	var repoKeys []string
	repoList := make(map[string]*data.Repository)
	for name, id := range author.Repositories {
		repo, err := data.GetRepository(ctx, dag, id)
		if err != nil {
			return nil, err
		}

		repoList[name] = repo
		repoKeys = append(repoKeys, name)
	}

	var followKeys []string
	followList := make(map[string]*data.Author)
	for _, id := range author.Following {
		author, err := s.node.Authors().Get(ctx, id)
		if err != nil && err != routing.ErrNotFound {
			return nil, err
		}

		followList[id.Pretty()] = author
		followKeys = append(followKeys, id.Pretty())
	}

	sort.Strings(repoKeys)
	sort.Strings(followKeys)

	var isFollowing bool
	for _, id := range cfg.Author.Following {
		isFollowing = isFollowing || id.Pretty() == peerID
	}

	return &ViewModel{
		Name: "author.html",
		Data: map[string]interface{}{
			"Author":      author,
			"IsFollowing": isFollowing,
			"FollowList":  followList,
			"FollowKeys":  followKeys,
			"RepoList":    repoList,
			"RepoKeys":    repoKeys,
			"PeerID":      peerID,
			"SelfID":      selfID,
		},
	}, nil
}

func (s *Author) Follow(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()
	cfg := s.node.Config()

	params := httprouter.ParamsFromContext(ctx)
	peerID := params.ByName("peer_id")

	pid, err := peer.Decode(peerID)
	if err != nil {
		return nil, err
	}

	set := make(map[string]peer.ID)
	for _, id := range cfg.Author.Following {
		set[id.String()] = id
	}
	set[pid.String()] = pid

	var following []peer.ID
	for _, id := range set {
		following = append(following, id)
	}

	cfg.Sequence++
	cfg.Author.Following = following

	if err := s.node.Authors().Subscribe(pid); err != nil {
		return nil, err
	}

	if err := cfg.Save(); err != nil {
		return nil, err
	}

	http.Redirect(w, req, "/"+peerID, http.StatusSeeOther)
	return nil, s.node.Authors().Publish(ctx)
}

func (s *Author) Unfollow(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()
	cfg := s.node.Config()

	params := httprouter.ParamsFromContext(ctx)
	peerID := params.ByName("peer_id")

	pid, err := peer.Decode(peerID)
	if err != nil {
		return nil, err
	}

	set := make(map[string]peer.ID)
	for _, id := range cfg.Author.Following {
		set[id.String()] = id
	}
	delete(set, pid.String())

	var following []peer.ID
	for _, id := range set {
		following = append(following, id)
	}

	cfg.Sequence++
	cfg.Author.Following = following

	if _, err := s.node.Authors().Unsubscribe(pid); err != nil {
		return nil, err
	}

	if err := cfg.Save(); err != nil {
		return nil, err
	}

	http.Redirect(w, req, "/"+peerID, http.StatusSeeOther)
	return nil, s.node.Authors().Publish(ctx)
}
