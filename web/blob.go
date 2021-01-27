package web

import (
	"html/template"
	"net/http"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-path"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

var blobView = template.Must(template.New("index.html").Funcs(funcs).ParseFS(templates, "templates/index.html", "templates/repo.html", "templates/_blob.html"))

type blobModel struct {
	ID     cid.Cid
	Branch string
	Blob   string
	Page   string
	Path   string
	Repo   *data.Repository
	Ref    string
	Tag    string
}

// Blob renders the repo blob view.
func (s *Server) Blob(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()
	params := httprouter.ParamsFromContext(ctx)

	name := params.ByName("name")
	file := params.ByName("file")
	ref := params.ByName("ref")

	id, err := s.store.GetCid(name)
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

	blob, err := unixfs.Cat(ctx, s.client, fnode.Cid())
	if err != nil {
		return err
	}

	model := blobModel{
		ID:   id,
		Blob: blob,
		Page: "blob",
		Path: file,
		Repo: repo,
		Ref:  ref,
	}

	if _, ok := repo.Branches[ref]; ok {
		model.Branch = ref
	}

	if _, ok := repo.Tags[ref]; ok {
		model.Tag = ref
	}

	return blobView.Execute(w, &model)
}
