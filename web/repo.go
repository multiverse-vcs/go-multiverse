package web

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-path"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

var (
	blobView    = template.Must(template.New("index.html").Funcs(funcs).ParseFS(templates, "templates/index.html", "templates/repo.html", "templates/_blob.html"))
	commitsView = template.Must(template.New("index.html").Funcs(funcs).ParseFS(templates, "templates/index.html", "templates/repo.html", "templates/_commits.html"))
	treeView    = template.Must(template.New("index.html").Funcs(funcs).ParseFS(templates, "templates/index.html", "templates/repo.html", "templates/_tree.html"))
)

var readmeRegex = regexp.MustCompile(`(?i)^readme.*`)

type repoModel struct {
	ID   cid.Cid
	Page string
	Path string
	Repo *data.Repository

	Blob    *blobModel
	Commits *commitsModel
	Tree    *treeModel
}

type blobModel struct {
	Data string
}

type commitsModel struct {
	IDs  []cid.Cid
	List []*data.Commit
}

type treeModel struct {
	Readme  string
	Entries []*unixfs.DirEntry
}

func (s *Server) Repo(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	parts := strings.Split(req.URL.Path, "/")

	for _, p := range parts {
		fmt.Println(p)
	}

	page := ""
	name := ""
	file := ""
	ref := ""

	if page == "" {
		return errors.New("invalid page")
	}

	id, err := s.store.GetCid(name)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return err
	}

	model := repoModel{
		ID:   id,
		Page: page,
		Path: file,
		Repo: repo,
	}

	switch page {
	case "blob":
		blob, err := s.Blob(ctx, repo, ref, file)
		if err != nil {
			return err
		}

		model.Blob = blob
		return blobView.Execute(w, &model)
	case "commits":
		commits, err := s.Commits(ctx, repo, ref)
		if err != nil {
			return err
		}

		model.Commits = commits
		return commitsView.Execute(w, &model)
	case "tree":
		tree, err := s.Tree(ctx, repo, ref, file)
		if err != nil {
			return err
		}

		model.Tree = tree
		return treeView.Execute(w, &model)
	}

	return errors.New("invalid page")
}

func (s *Server) Blob(ctx context.Context, repo *data.Repository, ref string, file string) (*blobModel, error) {
	var head cid.Cid

	fpath, err := path.FromSegments("/ipfs/", head.String(), "tree", file)
	if err != nil {
		return nil, err
	}

	fnode, err := s.client.ResolvePath(ctx, fpath)
	if err != nil {
		return nil, err
	}

	blob, err := unixfs.Cat(ctx, s.client, fnode.Cid())
	if err != nil {
		return nil, err
	}

	return &blobModel{
		Data: blob,
	}, nil
}

func (s *Server) Commits(ctx context.Context, repo *data.Repository, ref string) (*commitsModel, error) {
	var head cid.Cid

	var ids []cid.Cid
	visit := func(id cid.Cid) bool {
		ids = append(ids, id)
		return true
	}

	if err := core.Walk(ctx, s.client, head, visit); err != nil {
		return nil, err
	}

	var list []*data.Commit
	for _, id := range ids {
		commit, err := data.GetCommit(ctx, s.client, id)
		if err != nil {
			return nil, err
		}

		list = append(list, commit)
	}

	return &commitsModel{
		IDs:  ids,
		List: list,
	}, nil
}

func (s *Server) Tree(ctx context.Context, repo *data.Repository, ref string, file string) (*treeModel, error) {
	var head cid.Cid

	fpath, err := path.FromSegments("/ipfs/", head.String(), "tree", file)
	if err != nil {
		return nil, err
	}

	fnode, err := s.client.ResolvePath(ctx, fpath)
	if err != nil {
		return nil, err
	}

	entries, err := unixfs.Ls(ctx, s.client, fnode.Cid())
	if err != nil {
		return nil, err
	}

	readme, err := s.Readme(ctx, entries)
	if err != nil {
		return nil, err
	}

	return &treeModel{
		Readme:  readme,
		Entries: entries,
	}, nil
}

// Readme returns the contents of the readme if it exists.
func (s *Server) Readme(ctx context.Context, entries []*unixfs.DirEntry) (string, error) {
	for _, e := range entries {
		if readmeRegex.MatchString(e.Name) {
			return unixfs.Cat(ctx, s.client, e.Cid)
		}
	}
	return "", nil
}
