package view

import (
	"html/template"
	"net/http"

	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/node"
)

var homeView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("web/html/index.html", "web/html/home.html"))

type homeModel struct {
	List []*data.Repository
	node *node.Node
}

// Home returns the home view route.
func Home(node *node.Node) http.Handler {
	model := &homeModel{
		node: node,
	}

	return View(model.execute)
}

// execute renders the template as the http response.
func (model homeModel) execute(w http.ResponseWriter, req *http.Request) error {
	list, err := model.node.ListRepositories(req.Context())
	if err != nil {
		return err
	}

	model.List = list
	return homeView.Execute(w, &model)
}
