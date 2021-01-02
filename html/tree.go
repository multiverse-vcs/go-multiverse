package html

import (
	"context"
	"html/template"
	"io"
	gopath "path"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-path"
	"github.com/multiverse-vcs/go-multiverse/node"
)

var treeFunc = template.FuncMap{
	"join": func(parts ...string) string {
		return gopath.Join(parts...)
	},
}

var treeView = template.Must(template.New("index.html").Funcs(treeFunc).ParseFiles("html/index.html", "html/tree.html"))

type treeModel struct {
	Links []*ipld.Link
	Path  string
	Repo  string
	Util  *util
}

func Tree(ctx context.Context, w io.Writer, node *node.Node, repo, file string) error {
	id, err := node.Repo().Get(repo)
	if err != nil {
		return err
	}

	p, err := path.FromSegments("/ipfs/", id.String(), "tree", file)
	if err != nil {
		return err
	}

	tree, err := node.ResolvePath(ctx, p)
	if err != nil {
		return err
	}

	model := treeModel{
		Links: tree.Links(),
		Path:  file,
		Repo:  repo,
		Util:  &util{node},
	}

	return treeView.Execute(w, &model)
}
