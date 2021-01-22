package web

import (
	"context"
	"html/template"
	"net/http"
	"regexp"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-path"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/peer"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

var treeView = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("templates/index.html", "templates/repo.html", "templates/_tree.html"))

type treeModel struct {
	ID     cid.Cid
	Branch string
	Page   string
	Path   string
	Readme string
	Repo   *data.Repository
	Ref    string
	Tag    string
	Tree   []*unixfs.DirEntry
}

var readmeRegex = regexp.MustCompile(`(?i)^readme.*`)

// Readme returns the contents of the readme if it exists.
func Readme(ctx context.Context, client *peer.Client, tree []*unixfs.DirEntry) (string, error) {
	for _, e := range tree {
		if readmeRegex.MatchString(e.Name) {
			return unixfs.Cat(ctx, client, e.Cid)
		}
	}
	return "", nil
}

// Tree renders the repo tree view.
func (s *Server) Tree(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()
	params := httprouter.ParamsFromContext(ctx)

	sid := params.ByName("id")
	file := params.ByName("file")

	ref := params.ByName("ref")
	if ref == "" {
		ref = "default"
	}

	id, err := cid.Decode(sid)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return err
	}

	head, err := repo.Ref(ref)
	if err != nil {
		return err
	}

	fpath, err := path.FromSegments("/ipfs/", head.String(), "tree", file)
	if err != nil {
		return err
	}

	fnode, err := s.client.ResolvePath(ctx, fpath)
	if err != nil {
		return err
	}

	tree, err := unixfs.Ls(ctx, s.client, fnode.Cid())
	if err != nil {
		return err
	}

	readme, err := Readme(ctx, s.client, tree)
	if err != nil {
		return err
	}

	model := treeModel{
		ID:     id,
		Page:   "tree",
		Path:   file,
		Repo:   repo,
		Ref:    ref,
		Readme: readme,
		Tree:   tree,
	}

	if _, ok := repo.Branches[ref]; ok {
		model.Branch = ref
	}

	if _, ok := repo.Tags[ref]; ok {
		model.Tag = ref
	}

	return treeView.Execute(w, &model)
}
