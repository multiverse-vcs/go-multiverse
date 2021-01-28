package web

import (
	"html/template"
	"net/http"

	"github.com/ipfs/go-cid"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
)

var commitsView = template.Must(template.New("index.html").Funcs(funcs).ParseFS(templates, "templates/index.html", "templates/repo.html", "templates/_commits.html"))

type commitsModel struct {
	ID     cid.Cid
	IDs    []cid.Cid
	List   []*data.Commit
	Branch string
	Page   string
	Path   string
	Repo   *data.Repository
	Ref    string
	Tag    string
}

func (s *Server) Commits(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()
	params := httprouter.ParamsFromContext(ctx)

	name := params.ByName("name")
	ref := params.ByName("ref")
	file := params.ByName("file")

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

	var ids []cid.Cid
	visit := func(id cid.Cid) bool {
		ids = append(ids, id)
		return true
	}

	if err := core.Walk(ctx, s.client, head, visit); err != nil {
		return err
	}

	var list []*data.Commit
	for _, id := range ids {
		commit, err := data.GetCommit(ctx, s.client, id)
		if err != nil {
			return err
		}

		list = append(list, commit)
	}

	model := commitsModel{
		ID:   id,
		IDs:  ids,
		List: list,
		Page: "commits",
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

	return commitsView.Execute(w, &model)
}
