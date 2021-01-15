package view

import (
	"html/template"
	"net/http"

	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/node"
)

var homeView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("web/html/index.html", "web/html/home.html"))

type homeController struct {
	node *node.Node
}

type homeModel struct {
	List []*data.Repository
}

// Home returns the home view route.
func Home(node *node.Node) http.Handler {
	c := &homeController{
		node: node,
	}

	return View(c.ServeHTTP)
}

// ServeHTTP renders the template as the http response.
func (c *homeController) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	list, err := c.node.ListRepositories(ctx)
	if err != nil {
		return err
	}

	return homeView.Execute(w, &homeModel{list})
}
