package view

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/node"
	repoView "github.com/multiverse-vcs/go-multiverse/web/view/repo"
)

var (
	codeView    = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("web/html/index.html", "web/html/repo.html", "web/html/repo/code.html"))
	commitsView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("web/html/index.html", "web/html/repo.html", "web/html/repo/commits.html"))
)

type repoController struct {
	node *node.Node
}

type repoModel struct {
	Branch string
	Path   string
	Page   string
	Repo   *data.Repository
	Ref    string
	Tag    string
	URL    string

	Code    *repoView.CodeModel
	Commits *repoView.CommitsModel
}

// Repo returns the code view.
func Repo(node *node.Node) http.Handler {
	c := &repoController{
		node: node,
	}

	return View(c.ServeHTTP)
}

// ServeHTTP renders the template as the http response.
func (c *repoController) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()
	params := httprouter.ParamsFromContext(ctx)

	name := params.ByName("name")
	file := params.ByName("file")

	ref := params.ByName("ref")
	if ref == "" {
		ref = "default"
	}

	page := params.ByName("page")
	if page == "" {
		page = "tree"
	}

	repo, err := c.node.GetRepository(ctx, name)
	if err != nil {
		return err
	}

	id, err := repo.Ref(ref)
	if err != nil {
		return err
	}

	model := repoModel{
		Page: page,
		Path: file,
		Repo: repo,
		Ref:  ref,
		URL:  req.URL.Path,
	}

	if _, ok := repo.Branches[ref]; ok {
		model.Branch = ref
	}

	if _, ok := repo.Tags[ref]; ok {
		model.Tag = ref
	}

	switch page {
	case "commits":
		commits, err := repoView.Commits(ctx, c.node, id)
		if err != nil {
			return err
		}

		model.Commits = commits
		return commitsView.Execute(w, &model)
	case "tree":
		code, err := repoView.Code(ctx, c.node, id, file)
		if err != nil {
			return err
		}

		model.Code = code
		return codeView.Execute(w, &model)
	}

	return errors.New("page not found")
}
