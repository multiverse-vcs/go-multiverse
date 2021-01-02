package http

import (
	"html/template"
	"net/http"

	"github.com/ipfs/go-datastore/query"
)

var homeView = template.Must(template.ParseFiles("templates/index.html", "templates/home.html"))

type homeModel struct {
	List []query.Entry
}

func (s *Server) home(w http.ResponseWriter, req *http.Request) {
	res, err := s.dstore.Query(query.Query{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	all, err := res.Rest()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := homeModel{
		List: all,
	}

	if err := homeView.Execute(w, &model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
