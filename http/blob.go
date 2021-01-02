package http

import (
	"html/template"
	"net/http"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/core"
)

var blobView = template.Must(template.ParseFiles("templates/index.html", "templates/blob.html"))

type blobModel struct {
	Data string
	Path string
	Repo string
	Util *Server
}

func (s *Server) blob(w http.ResponseWriter, req *http.Request) {
	params := httprouter.ParamsFromContext(req.Context())

	repo := params.ByName("repo")
	file := params.ByName("file")

	val, err := s.dstore.Get(datastore.NewKey(repo))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := cid.Cast(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := path.FromSegments("/ipfs/", id.String(), "tree", file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resolver := resolver.NewBasicResolver(s.dag)
	node, err := resolver.ResolvePath(req.Context(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := core.Cat(req.Context(), s.dag, node.Cid())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := blobModel{
		Data: data,
		Path: file,
		Repo: repo,
		Util: s,
	}

	if err := blobView.Execute(w, &model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
