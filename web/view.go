package web

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"regexp"

	"github.com/ipfs/go-path"
	"github.com/julienschmidt/httprouter"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

var readmePattern = regexp.MustCompile(`(?i)^readme.*`)

var layout = template.Must(template.New("index.html").Funcs(funcs).ParseFS(templates, "html/*"))

// View is an http handler that renders a view.
type View func(http.ResponseWriter, *http.Request) (*ViewModel, error)

// ViewModel contains the data for the template.
type ViewModel struct {
	Name string
	Data interface{}
}

// ServeHTTP handles http requests to a route.
func (v View) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	model, err := v(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var page bytes.Buffer
	if err := layout.ExecuteTemplate(&page, model.Name, model.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := template.HTML(page.String())
	if err := layout.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Author renders the repository list.
func (s *Server) Author(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()

	params := httprouter.ParamsFromContext(ctx)
	peerID := params.ByName("peer_id")

	pid, err := peer.Decode(peerID)
	if err != nil {
		return nil, err
	}

	author, err := s.client.Authors().Search(ctx, pid)
	if err != nil {
		return nil, err
	}

	var keys []string
	var list []*data.Repository

	for name, id := range author.Repositories {
		repo, err := data.GetRepository(ctx, s.client, id)
		if err != nil {
			return nil, err
		}

		keys = append(keys, name)
		list = append(list, repo)
	}

	return &ViewModel{
		Name: "author.html",
		Data: map[string]interface{}{
			"Keys":   keys,
			"List":   list,
			"PeerID": peerID,
		},
	}, nil
}

// Tree renders the blob and tree file viewer.
func (s *Server) Tree(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()

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

	author, err := s.client.Authors().Search(ctx, pid)
	if err != nil {
		return nil, err
	}

	repoID, ok := author.Repositories[name]
	if !ok {
		return nil, errors.New("repository does not exist")
	}

	repo, err := data.GetRepository(ctx, s.client, repoID)
	if err != nil {
		return nil, err
	}

	fpath, err := path.FromSegments("/ipfs/", repoID.String(), refs, head, "tree", file)
	if err != nil {
		return nil, err
	}

	fnode, err := s.client.ResolvePath(ctx, fpath)
	if err != nil {
		return nil, err
	}

	dir, err := unixfs.IsDir(ctx, s.client, fnode.Cid())
	if err != nil {
		return nil, err
	}

	var blob string
	var tree []*unixfs.DirEntry

	var readmeBlob string
	var readmeName string

	switch {
	case dir:
		tree, err = unixfs.Ls(ctx, s.client, fnode.Cid())
		if err != nil {
			return nil, err
		}

		entry, err := unixfs.Find(ctx, s.client, fnode.Cid(), readmePattern)
		if err != nil || entry == nil {
			break
		}

		readmeName = entry.Name
		readmeBlob, err = unixfs.Cat(ctx, s.client, entry.Cid)
		if err != nil {
			return nil, err
		}
	default:
		blob, err = unixfs.Cat(ctx, s.client, fnode.Cid())
		if err != nil {
			return nil, err
		}
	}

	return &ViewModel{
		Name: "tree.html",
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
