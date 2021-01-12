package view

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"regexp"

	"github.com/ipfs/go-path"
	"github.com/ipfs/go-unixfs"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/node"
)

//var repoView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("web/html/index.html", "web/html/repo.html"))

type repoModel struct {
	Blob   string
	Branch string
	Commit *data.Commit
	IsDir  bool
	Path   string
	Repo   *data.Repository
	Tree   []*core.DirEntry
	URL    string
	node   *node.Node
}

// Repo returns the repo view route.
func Repo(node *node.Node) http.Handler {
	model := &repoModel{
		node: node,
	}

	return View(model.execute)
}

// Readme returns the contents of the readme if it exists.
func (model repoModel) Readme() (string, error) {
	for _, e := range model.Tree {
		matched, err := regexp.MatchString(`(?i)^readme.*`, e.Name)
		if err != nil {
			return "", err
		}

		if matched {
			return core.Cat(context.Background(), model.node, e.Cid)
		}
	}

	return "", nil
}

// execute renders the template as the http response.
func (model repoModel) execute(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()
	params := httprouter.ParamsFromContext(ctx)

	name := params.ByName("repo")
	file := params.ByName("file")

	repo, err := model.node.GetRepository(ctx, name)
	if err != nil {
		return err
	}

	branch := req.URL.Query().Get("branch")
	if branch == "" {
		branch = "default"
	}

	id, ok := repo.Branches[branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	c, err := model.node.Get(ctx, id)
	if err != nil {
		return err
	}

	commit, err := data.CommitFromCBOR(c.RawData())
	if err != nil {
		return err
	}

	p, err := path.FromSegments("/ipfs/", id.String(), "tree", file)
	if err != nil {
		return err
	}

	f, err := model.node.ResolvePath(ctx, p)
	if err != nil {
		return err
	}

	fsnode, err := unixfs.ExtractFSNode(f)
	if err != nil {
		return err
	}

	switch {
	case fsnode.IsDir():
		tree, err := core.Ls(ctx, model.node, f.Cid())
		if err != nil {
			return err
		}

		model.Tree = tree
	default:
		blob, err := core.Cat(ctx, model.node, f.Cid())
		if err != nil {
			return err
		}

		model.Blob = blob
	}

	model.Branch = branch
	model.Commit = commit
	model.IsDir = fsnode.IsDir()
	model.Path = file
	model.Repo = repo
	model.URL = req.URL.Path

	repoView := template.Must(template.New("index.html").Funcs(funcs).ParseFiles("web/html/index.html", "web/html/repo.html"))
	return repoView.Execute(w, &model)
}
