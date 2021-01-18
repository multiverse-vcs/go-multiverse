package view

import (
	"html/template"
	"net/http"

	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/node"
)

var homeView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("web/html/index.html", "web/html/home.html"))

type homeController struct {
	node  *node.Node
	store *data.Store
}

type homeModel struct {
	List []*data.Repository
}

// Home returns the home view route.
func Home(node *node.Node, store *data.Store) http.Handler {
	c := &homeController{
		node:  node,
		store: store,
	}

	return View(c.ServeHTTP)
}

// ServeHTTP renders the template as the http response.
func (c *homeController) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	keys, err := c.store.Keys()
	if err != nil {
		return err
	}

	var list []*data.Repository
	for _, key := range keys {
		id, err := c.store.GetCid(key)
		if err != nil {
			return err
		}

		repo, err := data.GetRepository(ctx, c.node, id)
		if err != nil {
			return err
		}

		list = append(list, repo)
	}

	return homeView.Execute(w, &homeModel{list})
}
