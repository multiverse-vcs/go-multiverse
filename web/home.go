package web

import (
	"html/template"
	"net/http"

	"github.com/multiverse-vcs/go-multiverse/data"
)

var homeView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("templates/index.html", "templates/home.html"))

type homeModel struct {
	Keys []string
	List []*data.Repository
}

// Home renders the home view.
func (s *Server) Home(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	keys, err := s.store.Keys()
	if err != nil {
		return err
	}

	var list []*data.Repository
	for _, k := range keys {
		id, err := s.store.GetCid(k)
		if err != nil {
			return err
		}

		repo, err := data.GetRepository(ctx, s.client, id)
		if err != nil {
			return err
		}

		list = append(list, repo)
	}

	model := homeModel{
		Keys: keys,
		List: list,
	}

	return homeView.Execute(w, &model)
}
