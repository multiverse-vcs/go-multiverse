package web

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"regexp"

	"github.com/ipfs/go-path"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

var readmePattern = regexp.MustCompile(`(?i)^readme.*`)

var layout = template.Must(template.New("index.html").Funcs(funcs).ParseFS(templates, "html/*"))

// View is an http handler that renders a view.
type View func(http.ResponseWriter, *http.Request) (*ViewModel, error)

// ViewModel contains the data for the template.
type ViewModel struct {
	Name string
	Page template.HTML
	Data interface{}
}

// ServeHTTP handles http requests to a route.
func (v View) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// uncomment below for development page reloads
	layout = template.Must(template.New("index.html").Funcs(funcs).ParseGlob("web/html/*"))

	model, err := v(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var page bytes.Buffer
	if err := layout.ExecuteTemplate(&page, model.Name, model.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	model.Page = template.HTML(page.String())
	if err := layout.Execute(w, &model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Home renders the repository list.
func (s *Server) Home(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()
	cfg := s.client.Config()

	var keys []string
	var list []*data.Repository

	for name, id := range cfg.Author.Repositories {
		repo, err := data.GetRepository(ctx, s.client, id)
		if err != nil {
			return nil, err
		}

		keys = append(keys, name)
		list = append(list, repo)
	}

	metrics, err := s.client.GetMetrics()
	if err != nil {
		return nil, err
	}

	return &ViewModel{
		Name: "home.html",
		Data: map[string]interface{}{
			"Keys":    keys,
			"List":    list,
			"Metrics": metrics,
		},
	}, nil
}

// Tree renders the blob and tree file viewer.
func (s *Server) Tree(w http.ResponseWriter, req *http.Request) (*ViewModel, error) {
	ctx := req.Context()
	cfg := s.client.Config()

	params := httprouter.ParamsFromContext(ctx)
	name := params.ByName("name")
	refs := params.ByName("refs")
	head := params.ByName("head")
	file := params.ByName("file")

	id, ok := cfg.Author.Repositories[name]
	if !ok {
		return nil, errors.New("repository does not exist")
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return nil, err
	}

	fpath, err := path.FromSegments("/ipfs/", id.String(), refs, head, "tree", file)
	if err != nil {
		return nil, err
	}

	fnode, err := s.client.ResolvePath(ctx, fpath)
	if err != nil {
		return nil, err
	}

	dir, err := unixfs.IsDir(ctx, s.client, fnode.Cid())
	if err != nil {
		return nil, err
	}

	var blob string
	var tree []*unixfs.DirEntry

	var readmeBlob string
	var readmeName string

	switch {
	case dir:
		tree, err = unixfs.Ls(ctx, s.client, fnode.Cid())
		if err != nil {
			return nil, err
		}

		entry, err := unixfs.Find(ctx, s.client, fnode.Cid(), readmePattern)
		if err != nil {
			return nil, err
		}

		if entry == nil {
			break
		}

		readmeName = entry.Name
		readmeBlob, err = unixfs.Cat(ctx, s.client, entry.Cid)
		if err != nil {
			return nil, err
		}
	default:
		blob, err = unixfs.Cat(ctx, s.client, fnode.Cid())
		if err != nil {
			return nil, err
		}
	}

	return &ViewModel{
		Name: "tree.html",
		Data: map[string]interface{}{
			"ID":         id,
			"Repo":       repo,
			"Refs":       refs,
			"File":       file,
			"Head":       head,
			"Blob":       blob,
			"Tree":       tree,
			"ReadmeBlob": readmeBlob,
			"ReadmeName": readmeName,
		},
	}, nil
}
