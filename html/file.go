package html

import (
	"context"
	"net/http"
	"regexp"

	"github.com/ipfs/go-path"
	"github.com/ipfs/go-unixfs"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
)

type fileModel struct {
	Cid    string
	Blob   string
	Tree   []*core.DirEntry
	Path   string
	Readme string
	Repo   string
	URL    string
	Util   *util
}

// File renders the file page.
func File(w http.ResponseWriter, req *http.Request, node *node.Node) error {
	ctx := req.Context()
	params := httprouter.ParamsFromContext(ctx)

	repo := params.ByName("repo")
	file := params.ByName("file")

	id, err := node.Repo().Get(repo)
	if err != nil {
		return err
	}

	p, err := path.FromSegments("/ipfs/", id.String(), "tree", file)
	if err != nil {
		return err
	}

	f, err := node.ResolvePath(ctx, p)
	if err != nil {
		return err
	}

	model := fileModel{
		Cid:  id.String(),
		Path: file,
		Repo: repo,
		URL:  req.URL.Path,
		Util: &util{node},
	}

	fsnode, err := unixfs.ExtractFSNode(f)
	if err != nil {
		return err
	}

	switch {
	case fsnode.IsDir():
		tree, err := core.Ls(ctx, node, f.Cid())
		if err != nil {
			return err
		}

		readme, err := readme(ctx, tree, node)
		if err != nil {
			return err
		}

		model.Tree = tree
		model.Readme = readme
		return compile("html/file_tree.html").Execute(w, &model)
	default:
		blob, err := core.Cat(ctx, node, f.Cid())
		if err != nil {
			return err
		}

		model.Blob = blob
		return compile("html/file_blob.html").Execute(w, &model)
	}
}

// readme returns the contents of the readme in the given directory tree if it exists.
func readme(ctx context.Context, tree []*core.DirEntry, node *node.Node) (string, error) {
	for _, e := range tree {
		matched, err := regexp.MatchString(`(?i)^readme\..*`, e.Name)
		if err != nil {
			return "", err
		}

		if matched {
			return core.Cat(ctx, node, e.Cid)
		}
	}

	return "", nil
}
