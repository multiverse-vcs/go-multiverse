package http

import (
	"html/template"
	"net/http"
	gopath "path"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
	"github.com/julienschmidt/httprouter"
)

var treeFunc = template.FuncMap{
	"join": func(parts ...string) string {
		return gopath.Join(parts...)
	},
}

var treeView = template.Must(template.New("index.html").Funcs(treeFunc).ParseFiles("templates/index.html", "templates/tree.html"))

type treeModel struct {
	Links []*ipld.Link
	Path  string
	Repo  string
	Util  *Server
}

func (s *Server) tree(w http.ResponseWriter, req *http.Request) {
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

	model := treeModel{
		Links: node.Links(),
		Path:  file,
		Repo:  repo,
		Util:  s,
	}

	if err := treeView.Execute(w, &model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
