package html

import (
	"net/http"

	"github.com/ipfs/go-datastore/query"
	"github.com/multiverse-vcs/go-multiverse/node"
)

type homeModel struct {
	List []query.Entry
}

// Home renders the home page.
func Home(w http.ResponseWriter, req *http.Request, node *node.Node) error {
	list, err := node.Repo().List()
	if err != nil {
		return err
	}

	model := homeModel{
		List: list,
	}

	return compile("html/home.html").Execute(w, &model)
}
