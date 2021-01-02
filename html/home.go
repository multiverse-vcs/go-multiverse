package html

import (
	"html/template"
	"io"

	"github.com/ipfs/go-datastore/query"
	"github.com/multiverse-vcs/go-multiverse/node"
)

var homeView = template.Must(template.ParseFiles("html/index.html", "html/home.html"))

type homeModel struct {
	List []query.Entry
}

// Home renders the home page.
func Home(w io.Writer, node *node.Node) error {
	list, err := node.Repo().List()
	if err != nil {
		return err
	}

	model := homeModel{
		List: list,
	}

	return homeView.Execute(w, &model)
}
