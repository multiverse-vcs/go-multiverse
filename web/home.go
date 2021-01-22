package web

import (
	"html/template"
	"net/http"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/data"
)

var homeView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("templates/index.html", "templates/home.html"))

type homeModel struct {
	IDs  []cid.Cid
	List []*data.Repository
}

// Home renders the home view.
func (s *Server) Home(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	pins, err := s.client.RecursiveKeys(ctx)
	if err != nil {
		return err
	}

	var list []*data.Repository
	for _, id := range pins {
		repo, err := data.GetRepository(ctx, s.client, id)
		if err != nil {
			return err
		}

		list = append(list, repo)
	}

	model := homeModel{
		IDs:  pins,
		List: list,
	}

	return homeView.Execute(w, &model)
}
